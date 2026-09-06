package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "mygrpc/user"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db map[string]*pb.UserResponse
}

// CreateUser (POST equivalent)
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	id := fmt.Sprintf("usr_%d", len(s.db)+1) // Simple ID generation

	newUser := &pb.UserResponse{
		Id:    id,
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}

	s.db[id] = newUser
	fmt.Printf("Server: Created user %s\n", id)
	return newUser, nil
}

// GetUser (GET equivalent)
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	u, exists := s.db[req.GetId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	fmt.Printf("Server: Found user %s\n", req.GetId())
	return u, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: make(map[string]*pb.UserResponse)})

	fmt.Println("Server running on port :50051...")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
