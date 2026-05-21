package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Sethuram52001/system-design-compendium/networking/exercises/grpc/solution/client-service/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
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

func (c *UserServiceClient) GetUser(id string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
