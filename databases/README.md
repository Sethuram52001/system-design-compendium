# Databases

## Overview

This section contains small, focused exercises for learning how database systems shape application design. The goal is not to build production-grade systems, but to understand the core tradeoffs behind storage models, access patterns, consistency, indexing, and operational behavior.

For now, this section is limited to Postgres.

Redis is intentionally kept in `caching/exercises/redis` instead of this folder. It can be used as a database in some systems, but for system design learning it is usually clearer to study it first as a cache, rate limiter, queue-like primitive, and coordination tool.

## Exercises

- [Postgres](exercises/postgres/README.md): relational modeling, constraints, transactions, indexes, locking, full-text search, geospatial indexes, and partitioning.

## Learning Objectives

- Understand how database choice affects system design.
- Practice selecting storage models based on access patterns.
- Learn when strong relational guarantees matter.
- Build intuition for read paths, write paths, indexes, locking, and operational constraints.

## Suggested Learning Path

1. Model a tiny Postgres domain with two normalized tables and constraints.
2. Add B-tree, composite, partial, expression, and unique indexes; inspect query behavior with `EXPLAIN`.
3. Practice transactions, row-level locking, optimistic locking, and pessimistic locking.
4. Add a small taste of full-text search, geospatial indexing, and partitioning.
