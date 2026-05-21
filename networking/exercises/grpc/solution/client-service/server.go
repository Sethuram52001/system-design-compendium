package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	grpcClient *UserServiceClient
}

func NewServer(grpcClient *UserServiceClient) *Server {
	return &Server{
		grpcClient: grpcClient,
	}
}

func (s *Server) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := s.grpcClient.GetUser(userID)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.GetId(),
		"name":  user.GetName(),
		"email": user.GetEmail(),
		"age":   user.GetAge(),
	})
}

func main() {
	grpcClient, err := NewUserServiceClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	server := NewServer(grpcClient)

	router := gin.Default()
	router.GET("/user/:id", server.GetUser)

	log.Println("ClientService HTTP server listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
