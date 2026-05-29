# Redis Exercise: Cache-Aside with In-Memory Data

## Exercise Description

Build a small Go HTTP service that uses Redis as a cache in front of an in-memory user store. This exercise keeps the source of truth intentionally simple so the focus stays on Redis behavior: cache hits, cache misses, TTLs, and invalidation.

You will not connect to Postgres or any other database in the main exercise.

## Main Exercise

Implement a tiny service with:

- An in-memory map as the source of truth.
- Redis as a cache layer.
- `GET /users/:id` with cache-aside behavior.
- `PATCH /users/:id` or `PUT /users/:id` that updates the in-memory store and invalidates the Redis key.
- Clear response metadata showing whether the request was a cache hit or miss.

## Tech Stack

- Go
- Redis
- `go-redis`
- Standard `net/http` or Gin
- In-memory map as the source of truth

## Requirements

- Start Redis locally.
- Connect the Go service to Redis.
- Seed a few users in memory.
- Implement `GET /users/:id`:
  - Check Redis first.
  - On hit, return cached data.
  - On miss, read from memory.
  - Store the user in Redis with a TTL.
- Implement update behavior:
  - Update the in-memory user.
  - Delete `user:{id}` from Redis.
- Return whether the response came from cache.
- Keep the main exercise focused on cache-aside and TTLs.

## Suggested Project Structure

```text
caching/exercises/redis/
├── README.md
├── IMPLEMENTATION_GUIDE.md
└── solution/
    ├── main.go
    ├── go.mod
    └── docker-compose.yml
```

## Learning Objectives

- Understand Redis as a cache, not the source of truth.
- Learn cache-aside reads.
- Learn TTL-based expiration.
- Practice cache invalidation after writes.
- Understand why stale data is a normal cache tradeoff.
- Build intuition for Redis keys and value serialization.

## Extensions

After the main exercise works, add one or more of these:

- **Sorted sets:** implement a leaderboard or top users by score using `ZADD`, `ZREVRANGE`, and `ZRANK`.
- **Rate limiting:** implement a tiny fixed-window limiter with `INCR` and `EXPIRE`.
- **Counters:** track endpoint request counts with `INCR`.
- **Pub/Sub with two services:** publish user update events from an API service and consume them in a separate notification or audit service.
- **Redis Streams in the same service:** append user update events with `XADD` and read recent events with `XREVRANGE`.
- **Cache metrics:** expose cache hit and miss counts.

Sorted sets are the recommended first extension because they show up frequently in system design discussions for leaderboards, ranking, trending items, priority queues, and time-windowed scoring.

Pub/Sub and Streams are intentionally extensions. Pub/Sub teaches live fanout between services. Streams teach persisted event history and consumer-style processing.
