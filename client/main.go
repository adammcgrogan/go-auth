package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go-auth/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Establish a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := auth.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run ./client <register|login> <username> <password>")
		os.Exit(1)
	}

	command := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]

	switch command {
	case "register":
		res, err := c.Register(ctx, &auth.RegisterRequest{Username: username, Password: password})
		if err != nil {
			log.Fatalf("could not register: %v", err)
		}
		fmt.Printf("User registered with ID: %s\n", res.UserId)
	case "login":
		res, err := c.Login(ctx, &auth.LoginRequest{Username: username, Password: password})
		if err != nil {
			log.Fatalf("could not login: %v", err)
		}
		fmt.Printf("Login successful. Token: %s\n", res.Token)
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
