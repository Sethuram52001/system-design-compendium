package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	CacheHit bool `json:"cache_hit"`
	User     User `json:"user"`
}

var (
	ctx         = context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	users = map[string]User{
		"1": {ID: "1", Name: "Alice", Email: "alice@gmail.com"},
		"2": {ID: "2", Name: "Bob", Email: "bob@gmail.com"},
	}
)

func userCacheKey(id string) string {
	return "user:" + id
}

func getUserFromStore(id string) (User, bool) {
	user, ok := users[id]
	return user, ok
}

func updateUserInStore(id string, user User) User {
	user.ID = id
	users[id] = user
	return user
}

func getCacheUser(id string) (User, bool, error) {
	cached, err := redisClient.Get(userCacheKey(id)).Result()
	if err == redis.Nil {
		return User{}, false, nil
	}

	if err != nil {
		return User{}, false, err
	}

	var user User
	if err := json.Unmarshal([]byte(cached), &user); err != nil {
		return User{}, false, err
	}

	return user, true, nil
}

func setCacheUser(user User) error {
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return redisClient.Set(userCacheKey(user.ID), payload, 5*time.Second).Err()
}

func invalidateUser(id string) error {
	return redisClient.Del(userCacheKey(id)).Err()
}

func getUserHandler(c *gin.Context) {
	id := c.Param("id")

	user, hit, err := getCacheUser(id)
	if err == nil && hit {
		c.JSON(200, UserResponse{CacheHit: true, User: user})
		return
	}

	user, ok := getUserFromStore(id)
	if !ok {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if err := setCacheUser(user); err != nil {
		log.Printf("redis set failed: %v", err)
	}

	c.JSON(200, UserResponse{CacheHit: false, User: user})
}

func updateUserHandler(c *gin.Context) {
	id := c.Param("id")

	var input User
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := updateUserInStore(id, input)

	if err := invalidateUser(user.ID); err != nil {
		log.Printf("redis delete failed: %v", err)
	}

	c.JSON(200, UserResponse{CacheHit: false, User: user})
}

func main() {
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	r := gin.Default()
	r.GET("/users/:id", getUserHandler)
	r.PUT("/users/:id", updateUserHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
