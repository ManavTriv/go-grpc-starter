package main

import (
	"context"
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