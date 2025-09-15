package main

import (
	"context"
	"fmt"
	"go-auth/auth"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtKey []byte

type server struct {
	auth.UnimplementedAuthServiceServer
	db *Database
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	log.Printf("[REGISTER] Received request for username: %s", req.Username)

	if len(req.Username) < 3 {
		return nil, status.Errorf(codes.InvalidArgument, "username must be at least 3 characters long")
	}
	if len(req.Password) < 8 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.db.CreateUser(req.Username, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &auth.RegisterResponse{UserId: fmt.Sprintf("%d", userID)}, nil
}

func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	log.Printf("[LOGIN] Received request for username: %s", req.Username)
	user, err := s.db.GetUserByUsername(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &auth.LoginResponse{Token: tokenString}, nil
}

// ListUsers implements the gRPC service endpoint for listing all registered users.
func (s *server) ListUsers(ctx context.Context, req *auth.ListUsersRequest) (*auth.ListUsersResponse, error) {
	log.Printf("[LISTUSERS] Received request to list all users")

	usernames, err := s.db.GetAllUsernames()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve users: %v", err)
	}

	return &auth.ListUsersResponse{Usernames: usernames}, nil
}

func main() {
	jwtSecret := getEnv("JWT_SECRET", "my_super_secret_default_key")
	jwtKey = []byte(jwtSecret)

	dbPath := getEnv("DB_PATH", "./users.db")
	port := getEnv("PORT", ":50051")

	db, err := NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, &server{db: db})

	go func() {
		log.Printf("Server listening on %s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
