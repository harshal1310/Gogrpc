package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mygrpc/user"
)

// CreateUserRequest makes a CreateUser RPC with the provided name and email.
func CreateUserRequest(name, email string) (*pb.UserResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	return res, nil
}

// GetUserRequest makes a GetUser RPC with the provided user ID.
func GetUserRequest(id string) (*pb.UserResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetUser(ctx, &pb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return res, nil
}
