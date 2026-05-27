# Redis Exercise: Implementation Guide

## Goal

Build a focused Redis cache-aside demo:

```text
HTTP request
  -> check Redis
  -> if miss, read in-memory map
  -> write Redis with TTL
  -> return user
```

The in-memory map is the source of truth. Redis is only a cache.

## Step 1: Start Redis

Create `docker-compose.yml` in `caching/exercises/redis/solution`:

```yaml
services:
  redis:
    image: redis:7
    ports:
      - "6379:6379"
```

Start Redis:

```bash
docker compose up
```

Optional sanity check:

```bash
redis-cli ping
```

Expected:

```text
PONG
```

## Step 2: Create the Go Service

From `caching/exercises/redis`:

```bash
mkdir -p solution
cd solution
go mod init github.com/Sethuram52001/system-design-compendium/caching/exercises/redis/solution
go get github.com/redis/go-redis/v9
go get github.com/gin-gonic/gin
```

Create `main.go`.

## Step 3: Add Imports and Models

Start with:

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
```

`User` is the domain object. `UserResponse` makes cache behavior visible to the caller.

## Step 4: Add the In-Memory Source of Truth

Use an in-memory map instead of a real database:

```go
var (
	users = map[string]User{
		"1": {ID: "1", Name: "Ada Lovelace", Email: "ada@example.com"},
		"2": {ID: "2", Name: "Grace Hopper", Email: "grace@example.com"},
	}
)
```

This map is intentionally unprotected because this is a small learning demo. A real service should protect shared in-memory state or use a real database.

Add small helper functions:

```go
func getUserFromStore(id string) (User, bool) {
	user, ok := users[id]
	return user, ok
}

func updateUserInStore(id string, user User) User {
	user.ID = id
	users[id] = user
	return user
}
```

## Step 5: Connect to Redis

Create a Redis client:

```go
var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)
```

In a production service, you would usually pass dependencies around instead of using package globals. For this learning exercise, globals keep the flow visible.

## Step 6: Write Cache Helpers

Use a predictable key format:

```go
func userCacheKey(id string) string {
	return "user:" + id
}
```

Read a cached user:

```go
func getCachedUser(id string) (User, bool, error) {
	cached, err := rdb.Get(ctx, userCacheKey(id)).Result()
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
```

Write a cached user with TTL:

```go
func setCachedUser(user User) error {
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, userCacheKey(user.ID), payload, 30*time.Second).Err()
}
```

Invalidate a cached user:

```go
func invalidateUser(id string) error {
	return rdb.Del(ctx, userCacheKey(id)).Err()
}
```

## Step 7: Implement Cache-Aside Read

For `GET /users/:id`, the handler should:

1. Extract the ID from the Gin route param.
2. Try Redis.
3. If Redis has the user, return `cache_hit: true`.
4. If Redis misses, read the in-memory map.
5. Store the user in Redis with TTL.
6. Return `cache_hit: false`.

Handler:

```go
func getUserHandler(c *gin.Context) {
	id := c.Param("id")

	user, hit, err := getCachedUser(id)
	if err == nil && hit {
		c.JSON(200, UserResponse{CacheHit: true, User: user})
		return
	}
	if err != nil {
		log.Printf("redis get failed: %v", err)
	}

	user, ok := getUserFromStore(id)
	if !ok {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	if err := setCachedUser(user); err != nil {
		log.Printf("redis set failed: %v", err)
	}

	c.JSON(200, UserResponse{CacheHit: false, User: user})
}
```

Notice the fallback behavior: if Redis has an error, this handler still tries the in-memory source of truth. That is a useful cache failure pattern.

## Step 8: Implement Update and Invalidation

For `PUT /users/:id`, update memory first, then delete the Redis key:

```go
func updateUserHandler(c *gin.Context) {
	id := c.Param("id")

	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid JSON"})
		return
	}

	updated := updateUserInStore(id, input)

	if err := invalidateUser(id); err != nil {
		log.Printf("redis delete failed: %v", err)
	}

	c.JSON(200, updated)
}
```

Why delete instead of update the cache directly?

For the first exercise, deletion is simpler and safer. The next read repopulates Redis from the source of truth.

## Step 9: Wire Routes

Gin removes most of the routing and JSON boilerplate:

```go
func main() {
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	router := gin.Default()
	router.GET("/users/:id", getUserHandler)
	router.PUT("/users/:id", updateUserHandler)

	log.Println("Redis cache-aside demo listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
```

## Step 10: Test the Main Flow

First request should miss:

```bash
curl http://localhost:8080/users/1
```

Expected shape:

```json
{
  "cache_hit": false,
  "user": {
    "id": "1",
    "name": "Ada Lovelace",
    "email": "ada@example.com"
  }
}
```

Second request should hit:

```bash
curl http://localhost:8080/users/1
```

Expected shape:

```json
{
  "cache_hit": true,
  "user": {
    "id": "1",
    "name": "Ada Lovelace",
    "email": "ada@example.com"
  }
}
```

Update the user:

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Byron","email":"ada.byron@example.com"}'
```

Then read again:

```bash
curl http://localhost:8080/users/1
```

Expected: `cache_hit` should be `false` because the update invalidated the old cache key.

To inspect Redis directly:

```bash
redis-cli GET user:1
redis-cli TTL user:1
```

## Extension: Sorted Sets

Sorted sets are one of Redis's most useful system design primitives. They store unique members with numeric scores.

Common uses:

- leaderboards
- top posts by score
- trending items
- priority queues
- recent activity scored by timestamp

Commands to learn:

```text
ZADD leaderboard:users 100 user:1
ZADD leaderboard:users 250 user:2
ZREVRANGE leaderboard:users 0 9 WITHSCORES
ZRANK leaderboard:users user:1
ZREVRANK leaderboard:users user:2
```

Add a score update endpoint:

```go
type ScoreRequest struct {
	Score float64 `json:"score"`
}

func updateScoreHandler(c *gin.Context) {
	id := c.Param("id")

	var input ScoreRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid JSON"})
		return
	}

	if _, ok := getUserFromStore(id); !ok {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	err := rdb.ZAdd(ctx, "leaderboard:users", redis.Z{
		Score:  input.Score,
		Member: "user:" + id,
	}).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update score"})
		return
	}

	c.JSON(200, gin.H{"user_id": id, "score": input.Score})
}
```

Read the top users:

```go
func leaderboardHandler(c *gin.Context) {
	top, err := rdb.ZRevRangeWithScores(ctx, "leaderboard:users", 0, 9).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read leaderboard"})
		return
	}

	c.JSON(200, top)
}
```

Suggested route wiring:

```go
router.POST("/users/:id/score", updateScoreHandler)
router.GET("/leaderboard", leaderboardHandler)
```

Gin can route `/users/:id/score` and `/users/:id` separately, so no manual path parsing is needed.

## Extension: Counters

Use `INCR` for request counts:

```go
rdb.Incr(ctx, "metrics:requests:/users")
```

A simple place to call this is at the top of `usersHandler`.

## Extension: Fixed-Window Rate Limiting

Use `INCR` plus `EXPIRE`:

```go
func allowRequest(clientID string, limit int64) bool {
	key := "rate_limit:" + clientID
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return true
	}
	if count == 1 {
		rdb.Expire(ctx, key, time.Minute)
	}
	return count <= limit
}
```

If `allowRequest` returns false, respond with HTTP `429 Too Many Requests`.

## Extension: Pub/Sub with Two Services

Use Pub/Sub to model live event fanout between processes.

Design:

```text
api-service
  -> updates in-memory user
  -> invalidates Redis cache
  -> publishes user.updated

notification-service / audit-service
  -> subscribes to user.updated
  -> logs or prints the update
```

Publisher snippet in the API service:

```go
func publishUserUpdated(id string) {
	event := map[string]any{
		"type":    "user.updated",
		"user_id": id,
	}
	payload, _ := json.Marshal(event)
	rdb.Publish(ctx, "users.events", payload)
}
```

Subscriber process:

```go
func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	sub := rdb.Subscribe(context.Background(), "users.events")
	defer sub.Close()

	for msg := range sub.Channel() {
		log.Printf("received event: %s", msg.Payload)
	}
}
```

Pub/Sub messages are not persisted. If the subscriber is offline, it misses the event.

## Extension: Redis Streams in the Same Service

Streams are better when you want an event history that can be read later.

On update:

```go
func appendUserEvent(id string) {
	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "user_events",
		Values: map[string]any{
			"type":    "user.updated",
			"user_id": id,
		},
	})
}
```

Read recent events:

```go
func eventsHandler(c *gin.Context) {
	events, err := rdb.XRevRangeN(ctx, "user_events", "+", "-", 10).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read events"})
		return
	}

	c.JSON(200, events)
}
```

Difference:

| Feature | Pub/Sub | Streams |
|---|---|---|
| Delivery | Live only | Stored in Redis |
| Read later | No | Yes |
| Multiple consumers | Basic fanout | Consumer groups available |
| Good for | live notifications | event logs, background processing |

For this exercise, keep Streams simple: append events and read recent events. Consumer groups can wait.

## Done Criteria

- Redis is running locally.
- `GET /users/:id` supports cache-aside reads.
- Cached values have a TTL.
- Updates invalidate the Redis key.
- Responses clearly show cache hit vs miss.
- You can explain why Redis is not the source of truth in this exercise.
- You have tried the sorted set extension or can explain how it would work.
