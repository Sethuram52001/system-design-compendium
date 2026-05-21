package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Sethuram52001/system-design-compendium/networking/exercises/grpc/solution/user-service/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	store *UserStore
}

func NewServer(store *UserStore) *Server {
	return &Server{store: store}
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	user, err := s.store.GetUser(req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{User: user}, nil
}

func main() {
	store := NewUserStore()

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, NewServer(store))

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("UserService gRPC server listening on :55051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
