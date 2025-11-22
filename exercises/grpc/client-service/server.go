package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server that handles HTTP requets and forwards them to the gRPC
type Server struct {
	grpcClient *UserServiceClient
}

// NewServer creates a new HTTP server
func NewServer(grpcClient *UserServiceClient) *Server {
	return &Server{
		grpcClient: grpcClient,
	}
}

// GetUser handles GET /user/:id requests
func (s *Server) GetUser(c *gin.Context) {
	userID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	log.Printf("HTTP request for user ID: %s", userID)

	// call gRPC service
	user, err := s.grpcClient.GetUser(userID)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})

			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid user ID",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}

		return
	}

	// return JSON response
	c.JSON(http.StatusOK, gin.H{
		"id":    user.GetId(),
		"name":  user.GetName(),
		"email": user.GetEmail(),
		"age":   user.GetAge(),
	})
}

func main() {
	// Connect to UserService gRPC server
	grpcClient, err := NewUserServiceClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	// Create HTTP server
	server := NewServer(grpcClient)

	// Create Gin router
	router := gin.Default()

	// Register routes
	router.GET("/user/:id", server.GetUser)

	// Start HTTP server
	log.Println("ClientService HTTP server listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
