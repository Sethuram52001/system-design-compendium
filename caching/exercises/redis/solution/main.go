package main

import (
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

type ScoreRequest struct {
	Score float64 `json:"score"`
}

var (
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
	countRequest("/users/:id")
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

	publishUserUpdate(user.ID)
	appendUserEvent(user.ID)

	c.JSON(200, UserResponse{CacheHit: false, User: user})
}

func updateScoreHandler(c *gin.Context) {
	id := c.Param("id")

	var input ScoreRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if _, ok := getUserFromStore(id); !ok {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	err := redisClient.ZAdd("leaderboard:users", redis.Z{
		Score:  input.Score,
		Member: "user:" + id,
	}).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update score"})
		return
	}

	c.JSON(200, gin.H{"user_id": id, "score": input.Score})

}

func leaderboardHandler(c *gin.Context) {
	top, err := redisClient.ZRevRangeWithScores("leaderboard:users", 0, 9).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch leaderboard"})
		return
	}

	c.JSON(200, gin.H{"leaderboard": top})
}

func alloweRequest(clientID string, limit int64) bool {
	key := "rate_limit:" + clientID

	count, err := redisClient.Incr(key).Result()
	if err != nil {
		log.Printf("redis incr failed: %v", err)
		return false
	}

	if count == 1 {
		redisClient.Expire(key, time.Minute)
	}

	return count <= limit
}

func rateLimitMiddleware(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.ClientIP()
		if !alloweRequest(clientID, limit) {
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func countRequest(path string) {
	redisClient.Incr("metrics:requests:" + path)
}

func publishUserUpdate(id string) {
	event := map[string]any{
		"type":    "user.updated",
		"user_id": id,
	}

	payload, _ := json.Marshal(event)
	redisClient.Publish("user.events", payload)
}

func appendUserEvent(id string) {
	redisClient.XAdd(&redis.XAddArgs{
		Stream: "user_events",
		Values: map[string]interface{}{
			"type":    "user.updated",
			"user_id": id,
		},
	})
}

func eventsHandler(c *gin.Context) {
	events, err := redisClient.XRevRangeN("user_events", "+", "-", 10).Result()

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch events"})
		return
	}
	c.JSON(200, gin.H{"events": events})
}

func main() {
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	r := gin.Default()
	r.Use(rateLimitMiddleware(5))

	r.GET("/users/:id", getUserHandler)
	r.PUT("/users/:id", updateUserHandler)

	// leaderboard
	r.GET("/leaderboard", leaderboardHandler)
	r.POST("/users/:id/scores", updateScoreHandler)

	// events
	r.GET("/events", eventsHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
