# Caching

## Overview

This section contains small exercises for learning cache and low-latency data access patterns. These systems often sit beside a primary database and are used to improve latency, reduce load, coordinate short-lived state, or implement high-throughput primitives.

Redis starts here because it is most useful to study first as a cache and ephemeral data-structure server. Later exercises can still explore Redis persistence and database-like use cases.

## Exercises

- [Redis](exercises/redis/README.md): cache-aside with an in-memory source of truth, TTLs, invalidation, and extensions for sorted sets and rate limiting.

## Wiki

- [Redis](wiki/redis.md)

## Resources

- [Caching resources](resources.md)

## Learning Objectives

- Understand why caches are introduced into system designs.
- Learn cache invalidation and expiration tradeoffs.
- Practice modeling short-lived data with TTLs.
- Explore rate limiting and counters.
- Understand where Redis helps and where a primary database is still required.
