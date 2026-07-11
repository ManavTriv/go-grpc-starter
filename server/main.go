package main

import (
	"context"
	"log"
	"net"
	"fmt"

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