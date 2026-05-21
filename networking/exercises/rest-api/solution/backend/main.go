package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User represents the user data structure
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// inMemoryStore stores users in memory
var (
	inMemoryStore = map[string]User{
		"1": {ID: "1", Name: "John Doe", Email: "john@example.com"},
		"2": {ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
	}
	nextID = 3
)

// getAllUsers handles GET /users - Retrieve all users
func getAllUsers(c *gin.Context) {
	var users []User
	for _, user := range inMemoryStore {
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// getUserByID handles GET /users/:id - Retrieve a user by ID
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	user, exists := inMemoryStore[id]

	if exists {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

// createUser handles POST /users - Create a new user
func createUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-generate ID if not provided
	if newUser.ID == "" {
		newUser.ID = strconv.Itoa(nextID)
		nextID++
	}

	// Check if user already exists
	if _, exists := inMemoryStore[newUser.ID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this ID already exists"})
		return
	}

	inMemoryStore[newUser.ID] = newUser
	c.JSON(http.StatusCreated, newUser)
}

// updateUser handles PUT /users/:id - Replace entire user data
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := inMemoryStore[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Ensure ID matches the URL parameter
	updatedUser.ID = id
	inMemoryStore[id] = updatedUser
	c.JSON(http.StatusOK, updatedUser)
}

// patchUser handles PATCH /users/:id - Partially update user data
func patchUser(c *gin.Context) {
	id := c.Param("id")

	user, exists := inMemoryStore[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}

	inMemoryStore[id] = user
	c.JSON(http.StatusOK, user)
}

// deleteUser handles DELETE /users/:id - Delete a user
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	if _, exists := inMemoryStore[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	delete(inMemoryStore, id)
	c.Status(http.StatusNoContent)
}

func main() {
	r := gin.Default()

	// Collection endpoints: /users
	r.GET("/users", getAllUsers)
	r.POST("/users", createUser)

	// Resource endpoints: /users/:id
	r.GET("/users/:id", getUserByID)
	r.PUT("/users/:id", updateUser)
	r.PATCH("/users/:id", patchUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
