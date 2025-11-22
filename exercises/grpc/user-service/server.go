package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Sethuram52001/system-design-compendium/exercises/grpc/user-service/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the UserService gRPC server
type Server struct {
	pb.UnimplementedUserServiceServer
	store *UserStore
}

// NewServer creates a new gRPC server instance
func NewServer(store *UserStore) *Server {
	return &Server{
		store: store,
	}
}

// GetUser implements the GetUser RPC method
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received GetUser request for ID: %s", req.GetId())

	// Validate request
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Retrieve user from store
	user, err := s.store.GetUser(req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func main() {
	// Create user store
	store := NewUserStore()

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register our service
	pb.RegisterUserServiceServer(grpcServer, NewServer(store))

	// Listen on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("UserService gRPC server listening on :50051")

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
