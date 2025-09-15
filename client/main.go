package main

import (
	"context"
	"fmt"
	"go-auth/auth"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run ./client <command> [args]")
	}

	// Establishes a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := auth.NewAuthServiceClient(conn)

	command := os.Args[1]

	// Use a switch statement to handle different client commands.
	switch command {
	case "register":
		if len(os.Args) != 4 {
			log.Fatalf("Usage: go run ./client register <username> <password>")
		}
		username := os.Args[2]
		password := os.Args[3]
		res, err := c.Register(context.Background(), &auth.RegisterRequest{Username: username, Password: password})
		if err != nil {
			log.Fatalf("could not register: %v", err)
		}
		log.Printf("User registered with ID: %s", res.UserId)

	case "login":
		if len(os.Args) != 4 {
			log.Fatalf("Usage: go run ./client login <username> <password>")
		}
		username := os.Args[2]
		password := os.Args[3]
		res, err := c.Login(context.Background(), &auth.LoginRequest{Username: username, Password: password})
		if err != nil {
			log.Fatalf("could not login: %v", err)
		}
		log.Printf("Login successful. Token: %s", res.Token)

	case "list":
		if len(os.Args) != 2 {
			log.Fatalf("Usage: go run ./client list")
		}
		res, err := c.ListUsers(context.Background(), &auth.ListUsersRequest{})
		if err != nil {
			log.Fatalf("could not list users: %v", err)
		}
		fmt.Println("Registered Users:")
		for _, username := range res.Usernames {
			fmt.Printf("- %s\n", username)
		}

	default:
		log.Fatalf("Unknown command: %s. Available commands: register, login, list", command)
	}
}
