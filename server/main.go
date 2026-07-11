package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "go-grpc-starter/gen/userpb"
)

type server struct {
	userpb.UnimplementedUserServiceServer
	users map[string]*userpb.User
}

func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user with id %s not found", req.Id)
	}
	return user, nil
}

func (s *server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	id := fmt.Sprintf("user-%d", len(s.users)+1)
	user := &userpb.User{
		Id: id,
		Name: req.Name,
		Email: req.Email,
	}
	s.users[id] = user
	return user, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	
	// Create our server instance, with an empty map ready to store users
	myServer := &server{
		users: make(map[string]*userpb.User),
	}
	// Tell the gRPC server to route UserService requests to myServer
	userpb.RegisterUserServiceServer(s, myServer)

	log.Println("Server listening on :50051")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}