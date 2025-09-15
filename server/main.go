package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os" // Import the 'os' package
	"time"

	"go-auth/auth"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

// Declare jwtKey at the package level so all functions can see it.
var jwtKey []byte

// server is used to implement auth.AuthServiceServer.
type server struct {
	auth.UnimplementedAuthServiceServer
	db *Database
}

// Claims struct for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Register creates a new user in the database.
func (s *server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	log.Printf("[REGISTER] Request for username: %s", req.Username)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[REGISTER] Error hashing password: %v", err)
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	userID, err := s.db.CreateUser(req.Username, string(hashedPassword))
	if err != nil {
		log.Printf("[REGISTER] Error creating user: %v", err)
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	log.Printf("[REGISTER] User registered with ID: %d", userID)
	return &auth.RegisterResponse{UserId: fmt.Sprintf("%d", userID)}, nil
}

// Login authenticates a user and returns a JWT.
func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	log.Printf("[LOGIN] Request for username: %s", req.Username)

	user, err := s.db.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("[LOGIN] Invalid credentials: %v", err)
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("[LOGIN] Invalid credentials")
		return nil, fmt.Errorf("invalid credentials")
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
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	log.Printf("[LOGIN] User %s authenticated successfully", req.Username)
	return &auth.LoginResponse{Token: tokenString}, nil
}

// Helper function to get an environment variable with a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	jwtSecret := getEnv("JWT_SECRET", "my_super_secret_default_key")
	jwtKey = []byte(jwtSecret)

	dbPath := getEnv("DB_PATH", "./users.db")
	db, err := NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("[DB] failed to initialize database: %v", err)
	}
	defer db.Close()

	port := getEnv("PORT", ":50051")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, &server{db: db})

	log.Printf("Server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
