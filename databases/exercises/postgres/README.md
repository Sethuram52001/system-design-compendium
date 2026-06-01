# Postgres Exercise: Indexes, Transactions, and Locking

## Exercise Description

Build a SQL-first Postgres lab. The goal is to learn database behavior directly through SQL files, not through a CRUD API.

You will create the solution files yourself inside `solution/`. Use `IMPLEMENTATION_GUIDE.md` when you want hints and `temp.md` only as a temporary full reference.

## Tech Stack

- Postgres
- SQL files
- `psql`
- Docker Compose

## Starter Structure

```text
databases/exercises/postgres/
├── README.md
├── IMPLEMENTATION_GUIDE.md
├── temp.md
└── solution/
    └── docker-compose.yml
```

## Exercises

### Exercise 1: Core Schema

Create `solution/01_schema.sql`.

Requirements:

- Create a `users` table.
- Create an `orders` table.
- Keep the schema to these two tables.
- Add primary keys, a foreign key from `orders.user_id` to `users.id`, `NOT NULL`, `UNIQUE`, `CHECK`, and timestamp defaults.
- Add a small amount of seed data.

### Exercise 2: Transactions

Create `solution/02_transactions.sql`.

Requirements:

- Write one transaction that creates a user and an order together.
- Write one rollback example.
- Test that rolled-back data does not persist.
- Test a case where user creation fails and the order write does not happen.
- Test a case where user creation succeeds but order creation fails, causing the user insert to roll back too.

### Exercise 3: Basic Indexes

Create `solution/03_indexes.sql`.

Requirements:

- Run `EXPLAIN` or `EXPLAIN ANALYZE` before adding indexes.
- Add a B-tree index.
- Add a composite index.
- Add a partial index.
- Add an expression index.
- Observe the query plan again after adding indexes.
- Write a short note for yourself explaining which query shape each index helps.

### Exercise 4: Pessimistic Locking

Create:

- `solution/04_locking_terminal_a.sql`
- `solution/04_locking_terminal_b.sql`

Requirements:

- Use `SELECT ... FOR UPDATE` to lock a row in one transaction.
- In another terminal, try to update the same row.
- Observe that the second update waits until the first transaction releases the lock.

### Exercise 5: Optimistic Locking and `SKIP LOCKED`

Create `solution/05_optimistic_and_skip_locked.sql`.

Requirements:

- Use a `version` column for optimistic locking.
- Make an update succeed only when the expected version matches.
- Use `FOR UPDATE SKIP LOCKED` to simulate workers claiming pending orders.

### Exercise 6: Full-Text Search

Create `solution/06_full_text_search.sql`.

Requirements:

- Create a small `articles` table.
- Add a generated `tsvector` column.
- Add a GIN index.
- Query using `websearch_to_tsquery` or `to_tsquery`.

### Exercise 7: Geospatial Index

Create `solution/07_geospatial_postgis.sql`.

Requirements:

- Use PostGIS.
- Create a small `places` table.
- Store a point location.
- Add a GiST index.
- Query for places within a distance.

### Exercise 8: Partitioning and BRIN

Create `solution/08_partitioning_brin.sql`.

Requirements:

- Create an append-style `events` table partitioned by timestamp range.
- Create at least two range partitions.
- Insert sample events.
- Add a normal B-tree index on one partition.
- Add a BRIN index on another partition.
- Use `EXPLAIN` to observe partition pruning.

## Done Criteria

- You can run every SQL file you created from `solution/`.
- You can explain when to use B-tree, composite, partial, expression, GIN, GiST, and BRIN indexes.
- You have tested commit, rollback, row-level locking, optimistic locking, and `SKIP LOCKED`.
- You can explain why partitioning changes table layout but does not automatically make every query faster.
