# Redis

## What It Is

Redis is an in-memory data store commonly used for caching, counters, rate limiting, queues, leaderboards, sessions, and short-lived coordination data.

It is fast because most operations happen in memory and Redis provides simple data structures with efficient commands.

## Why It Exists

Many systems need very low-latency access to data that is either temporary, frequently read, expensive to recompute, or useful as a derived view.

Redis is often placed beside a primary database to reduce read load and improve latency. In that setup, the database remains the source of truth and Redis acts as a speed layer.

## Mental Model

Think of Redis as a fast data-structure server.

Instead of only storing opaque values, Redis gives you primitives:

| Type | Common Use |
|---|---|
| String | cached JSON, counters |
| Hash | object-like fields |
| List | queues, recent items |
| Set | uniqueness, membership |
| Sorted set | rankings, leaderboards, scored timelines |
| Pub/Sub | lightweight message fanout |
| Stream | append-only event log and consumer processing |

## Cache-Aside

Cache-aside is the most common first Redis pattern:

```text
read request
  -> check Redis
  -> if hit, return cached value
  -> if miss, read source of truth
  -> write value to Redis with TTL
  -> return value
```

On writes, update the source of truth first and invalidate the cache key.

This keeps Redis from becoming the primary database by accident.

## TTLs and Invalidation

TTL means time to live. It lets Redis automatically expire keys after a duration.

TTL helps because cached data can become stale. Invalidation helps because known writes should remove or refresh old cached values immediately.

Common strategy:

```text
read: populate cache with TTL
write: update source of truth, then delete cache key
```

## Sorted Sets

Sorted sets are a Redis structure where each member has a numeric score.

They are useful for:

- leaderboards
- trending posts
- top-N queries
- ranking users
- priority queues
- timestamp-scored activity feeds

Example commands:

```text
ZADD leaderboard 100 user:1
ZADD leaderboard 250 user:2
ZREVRANGE leaderboard 0 9 WITHSCORES
```

Sorted sets show up often in system design because they make ranking and top-K lookups simple and fast.

## Pub/Sub

Redis Pub/Sub sends messages to currently connected subscribers.

It is useful for live fanout:

```text
service A publishes event
service B receives event if subscribed right now
```

Pub/Sub is not durable. If a subscriber is offline, it misses the message. That makes it useful for live notifications and coordination, but not for jobs that must be processed eventually.

## Streams

Redis Streams store events in an append-only log-like structure.

They are useful when you want to write events and read them later:

```text
XADD user_events * type user.updated user_id 1
XREAD STREAMS user_events 0
```

Streams can also support consumer groups, which makes them closer to a lightweight queue or event processing system than Pub/Sub.

## Pub/Sub vs Streams

| Question | Pub/Sub | Streams |
|---|---|---|
| Must subscribers be online? | Yes | No |
| Are messages retained? | No | Yes |
| Can consumers replay history? | No | Yes |
| Best for | live fanout | event logs, background workers |
| Complexity | simpler | more powerful |

## When To Use Redis

Use Redis when:

- reads are frequent and data can tolerate short-lived staleness
- expensive values can be cached
- counters or rate limits need low latency
- ranking or top-N queries fit sorted sets
- live fanout can use Pub/Sub
- lightweight event history can use Streams
- short-lived session or coordination state is needed

## When Not To Use Redis

Avoid Redis as the default choice when:

- strong durability is required
- relational constraints or joins matter
- the data must never be stale
- the working set is too large for memory
- a simple database index would solve the problem

## Common Pitfalls

- Treating Redis as the source of truth without planning persistence and recovery.
- Forgetting TTLs and growing memory forever.
- Invalidating the wrong key or forgetting invalidation.
- Creating hot keys that every request reads or writes.
- Caching data that changes too frequently.
- Adding Redis before measuring whether the database is actually a bottleneck.

## Related Concepts

- Postgres
- Cache invalidation
- Rate limiting
- Leaderboards
- Pub/Sub
- Streams
- Hot keys

## Exercise

- [Redis Exercise](../exercises/redis/README.md)
