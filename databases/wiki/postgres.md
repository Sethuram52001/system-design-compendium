# Postgres

## What It Is

Postgres is a relational database commonly used as a system of record for application data.

It gives you durable storage, SQL queries, constraints, transactions, indexes, concurrency control, and extensions for specialized use cases like full-text search and geospatial lookup.

## Why It Exists

Most applications need a reliable place to store important state: users, orders, payments, accounts, permissions, inventory, and other business data.

Postgres is useful because it combines strong correctness guarantees with practical query power. You can model relationships, enforce data integrity, update multiple rows safely, and ask complex questions using SQL.

## Mental Model

Think of Postgres as a durable relational engine:

```text
tables hold structured data
constraints protect correctness
transactions protect multi-step changes
indexes speed up selected query patterns
locks coordinate concurrent writers
```

Your schema and indexes should come from access patterns. A table can store the data correctly, but the right index makes the common reads efficient.

## Tables and Constraints

Tables define the shape of stored data. Constraints define what valid data means.

Common constraints:

| Constraint | Purpose |
|---|---|
| `PRIMARY KEY` | uniquely identifies each row |
| `FOREIGN KEY` | protects relationships between tables |
| `NOT NULL` | requires a value |
| `UNIQUE` | prevents duplicates |
| `CHECK` | validates allowed values or ranges |
| `DEFAULT` | supplies a value when one is omitted |

Constraints matter because they protect the database even when application code has bugs or multiple services write to the same data.

## Transactions

A transaction groups multiple statements into one unit:

```text
BEGIN
  change A
  change B
COMMIT or ROLLBACK
```

`COMMIT` makes the changes permanent. `ROLLBACK` discards all changes since `BEGIN`.

Transactions prevent partial workflows. For example, if user creation succeeds but order creation fails, rollback can remove the user insert too.

## Indexes

Indexes are extra data structures that help Postgres find rows faster.

Common index types and patterns:

| Index | Useful For |
|---|---|
| B-tree | equality and range queries |
| Composite | queries filtering by multiple columns |
| Partial | indexing only a subset of rows |
| Expression | indexing computed expressions like `lower(email)` |
| GIN | inverted-index-style lookup for `tsvector`, arrays, JSONB |
| GiST | geospatial and range-like searches |
| BRIN | large append-heavy tables with naturally ordered data |

Indexes speed up reads but add storage and write overhead. Add indexes for query shapes you actually care about.

## Query Planning

`EXPLAIN` shows how Postgres plans to run a query. `EXPLAIN ANALYZE` runs the query and shows actual timing.

Use it to answer:

- Is Postgres scanning the whole table?
- Is it using an index?
- How many rows does it expect?
- How many rows did it actually process?

For tiny tables, Postgres may choose a sequential scan even when an index exists. That can be correct because reading the whole tiny table is cheaper.

## Locking and Concurrency

Postgres uses locks to keep concurrent transactions from corrupting data.

Important patterns:

- `SELECT ... FOR UPDATE` locks selected rows for update.
- Pessimistic locking blocks other writers until the first transaction finishes.
- Optimistic locking uses a `version` column and updates only if the version still matches.
- `FOR UPDATE SKIP LOCKED` lets workers skip rows already claimed by another transaction.

In system design, these patterns show up in inventory updates, payment workflows, job queues, order state transitions, and account balance changes.

## Full-Text Search

Postgres has built-in full-text search using `tsvector` and `tsquery`.

The usual shape is:

```text
text columns -> tsvector -> GIN index -> text search query
```

This is useful for small or moderate search needs when you do not yet need a dedicated search system.

## Geospatial Lookup

With PostGIS, Postgres can store and query locations.

Common use cases:

- find nearby drivers
- find stores within a radius
- search places inside a region
- rank results by distance

GiST indexes help make spatial queries efficient.

## Partitioning

Partitioning splits one logical table into smaller physical tables.

Time-range partitioning is common for event-style data:

```text
events
  events_2026_05
  events_2026_06
```

Partitioning can help with retention, maintenance, and pruning irrelevant data. It is not a magic performance switch; queries need filters that match the partition key.

## When To Use Postgres

Use Postgres when:

- relational modeling matters
- strong consistency matters
- constraints protect important data
- transactions are needed for multi-step workflows
- SQL queries and joins are useful
- one reliable primary database is enough to start

## When To Be Careful

Be careful when:

- writes are extremely high volume
- data is globally distributed with low-latency writes everywhere
- queries require specialized search ranking at large scale
- tables grow without archiving, partitioning, or retention plans
- indexes are added without understanding write overhead

## Common Pitfalls

- Skipping database constraints and relying only on application validation.
- Adding indexes without checking query patterns.
- Assuming every index will be used.
- Forgetting that indexes slow down writes.
- Holding transactions open too long.
- Using locks without thinking about contention and deadlocks.
- Partitioning before the data volume or retention problem is real.
- Assuming multi-region Postgres means multi-region writes. Standard Postgres replication is usually single-primary: one leader accepts writes and replicas serve reads or standby failover. Deploying replicas in multiple locations can improve read latency and availability, but it does not remove the single write-leader bottleneck without extra architecture or distributed Postgres-style systems.

## Related Concepts

- Transactions
- Indexes
- Query planning
- Locking
- Full-text search
- Geospatial indexing
- Partitioning
- Caching

## Exercise

- [Postgres Exercise](../exercises/postgres/README.md)
