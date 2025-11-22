package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Sethuram52001/system-design-compendium/exercises/grpc/client-service/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserServiceClient struct wraps the gRPC client
type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserServiceClient creates a new gRPC client connection
func NewUserServiceClient(address string) (*UserServiceClient, error) {
	//	connect to the gRPC server
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := pb.NewUserServiceClient(conn)

	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

// GetUser calls the UserService GetUser RPC
func (c *UserServiceClient) GetUser(id string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetUserRequest{
		Id: id,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
