package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "go-grpc-starter/gen/userpb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new user
	createReq := &userpb.CreateUserRequest{
		Name: "John Doe",
		Email: "john.doe@example.com",
	}
	newUser, err := client.CreateUser(ctx, createReq)
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf("Created user: %v", newUser)

	// Fetch the user we just created 
	getReq := &userpb.GetUserRequest{
		Id: newUser.Id,
	}
	fetchedUser, err := client.GetUser(ctx, getReq)
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("Fetched user: %v", fetchedUser)
}