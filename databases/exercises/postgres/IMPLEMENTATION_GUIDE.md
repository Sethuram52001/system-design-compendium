# Postgres Exercise: Implementation Guide

## Goal

Complete the Postgres exercises by creating your own SQL files inside `solution/`.

This guide gives the path, important snippets, and minimal run commands. It is intentionally not a full copy of every final file. Use `temp.md` if you want to compare against complete reference files.

## Basic Run Loop

Start Postgres:

```bash
cd databases/exercises/postgres/solution
docker compose up -d
```

Run any SQL file you create:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 01_schema.sql
```

If local `psql` is not installed:

```bash
docker compose exec -T postgres psql -U app -d system_design_postgres < 01_schema.sql
```

For most exercises, the loop is:

1. Create or edit the SQL file in `solution/`.
2. Run it with `psql -f`.
3. Read the output.
4. Adjust and rerun.

## Exercise 1: Core Schema

Create `solution/01_schema.sql`.

Start by dropping child tables before parent tables so reruns are easy:

```sql
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
```

Create `users` with identity, uniqueness, status validation, and timestamps:

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'disabled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Create `orders` with a foreign key to `users`:

```sql
CREATE TABLE orders (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  status TEXT NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'paid', 'cancelled')),
  total_cents INTEGER NOT NULL CHECK (total_cents >= 0),
  version INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Add a few `INSERT` statements so later exercises have data.

Run:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 01_schema.sql
```

## Exercise 2: Transactions

Create `solution/02_transactions.sql`.

Use a transaction when multiple statements must succeed or fail together:

```sql
BEGIN;

-- insert a user
-- insert an order for that user

COMMIT;
```

A common beginner trap is hardcoding the user ID returned by the first insert. Prefer `RETURNING` with a CTE:

```sql
WITH new_user AS (
  INSERT INTO users (email, name)
  VALUES ('linus@example.com', 'Linus')
  RETURNING id
)
INSERT INTO orders (user_id, status, total_cents)
SELECT id, 'pending', 4200
FROM new_user;
```

Add a rollback example:

```sql
BEGIN;

INSERT INTO users (email, name)
VALUES ('rollback@example.com', 'Rollback Demo');

ROLLBACK;
```

Verify the rollback:

```sql
SELECT *
FROM users
WHERE email = 'rollback@example.com';
```

Now add a failure case where the first write fails. For example, try inserting an email that already exists from your seed data:

```sql
BEGIN;

INSERT INTO users (email, name)
VALUES ('ada@example.com', 'Duplicate Ada');

INSERT INTO orders (user_id, status, total_cents)
VALUES (1, 'pending', 1000);

ROLLBACK;
```

The duplicate email violates the `UNIQUE` constraint. After that error, Postgres marks the transaction as failed, so the order insert does not successfully happen either.

Then add the more useful rollback case: the user insert succeeds, but the order insert fails. This proves the transaction protects earlier successful work from being partially committed:

```sql
BEGIN;

INSERT INTO users (email, name)
VALUES ('order-fails@example.com', 'Order Fails');

INSERT INTO orders (user_id, status, total_cents)
SELECT id, 'pending', -100
FROM users
WHERE email = 'order-fails@example.com';

ROLLBACK;

SELECT *
FROM users
WHERE email = 'order-fails@example.com';
```

The order fails because `total_cents` violates the `CHECK (total_cents >= 0)` constraint. The final `SELECT` should return no rows because the transaction rolled back the successful user insert too.

Run:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 02_transactions.sql
```

## Exercise 3: Basic Indexes

Create `solution/03_indexes.sql`.

First run a query plan before adding indexes:

```sql
EXPLAIN ANALYZE
SELECT *
FROM orders
WHERE user_id = 1
  AND status = 'pending';
```

Add indexes based on query shapes:

```sql
CREATE INDEX orders_user_id_idx ON orders(user_id);
CREATE INDEX orders_user_status_idx ON orders(user_id, status);
CREATE INDEX orders_pending_user_idx ON orders(user_id) WHERE status = 'pending';
CREATE INDEX users_lower_email_idx ON users (lower(email));
```

Run the same `EXPLAIN ANALYZE` queries again. With tiny seed data, Postgres may still choose a sequential scan. That is okay; the important part is learning which index matches which query shape.

Run:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 03_indexes.sql
```

## Exercise 4: Pessimistic Locking

Create `solution/04_locking_terminal_a.sql` and `solution/04_locking_terminal_b.sql`.

Terminal A should start a transaction and lock one row:

```sql
BEGIN;

SELECT *
FROM orders
WHERE id = 1
FOR UPDATE;

SELECT pg_sleep(20);

COMMIT;
```

Terminal B should try to update the same row:

```sql
UPDATE orders
SET status = 'paid'
WHERE id = 1;
```

Run terminal A first:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 04_locking_terminal_a.sql
```

While terminal A is sleeping, run terminal B:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 04_locking_terminal_b.sql
```

Terminal B should wait until terminal A commits.

## Exercise 5: Optimistic Locking and `SKIP LOCKED`

Create `solution/05_optimistic_and_skip_locked.sql`.

Optimistic locking checks a version instead of blocking first:

```sql
UPDATE orders
SET status = 'cancelled',
    version = version + 1
WHERE id = 2
  AND version = 1;
```

If the update affects zero rows, treat it as a conflict.

Use `SKIP LOCKED` to claim work without waiting on already locked rows:

```sql
BEGIN;

SELECT id
FROM orders
WHERE status = 'pending'
ORDER BY id
FOR UPDATE SKIP LOCKED
LIMIT 1;

COMMIT;
```

Run:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 05_optimistic_and_skip_locked.sql
```

## Exercise 6: Full-Text Search

Create `solution/06_full_text_search.sql`.

Use `tsvector` for searchable text and GIN for term lookup:

```sql
CREATE TABLE articles (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  search_vector tsvector GENERATED ALWAYS AS (
    to_tsvector('english', title || ' ' || body)
  ) STORED
);

CREATE INDEX articles_search_gin_idx
ON articles
USING GIN (search_vector);
```

Query it:

```sql
SELECT title
FROM articles
WHERE search_vector @@ websearch_to_tsquery('english', 'postgres indexes');
```

Run:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 06_full_text_search.sql
```

## Exercise 7: Geospatial Index

Create `solution/07_geospatial_postgis.sql`.

PostGIS is not available in the plain `postgres` image. For this exercise, change `docker-compose.yml` to use a PostGIS image such as:

```yaml
image: postgis/postgis:16-3.4
```

Use `geography(Point, 4326)` and a GiST index:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE places (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  location geography(Point, 4326) NOT NULL
);

CREATE INDEX places_location_gist_idx
ON places
USING GIST (location);
```

Query by distance:

```sql
SELECT name
FROM places
WHERE ST_DWithin(
  location,
  ST_MakePoint(-73.9857, 40.7484)::geography,
  5000
);
```

## Exercise 8: Partitioning and BRIN

Create `solution/08_partitioning_brin.sql`.

Partition an append-style table by time:

```sql
CREATE TABLE events (
  id BIGSERIAL,
  created_at TIMESTAMPTZ NOT NULL,
  payload JSONB NOT NULL
) PARTITION BY RANGE (created_at);
```

Create child partitions:

```sql
CREATE TABLE events_2026_05
PARTITION OF events
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
```

Use BRIN for naturally ordered, append-heavy timestamp data:

```sql
CREATE INDEX events_2026_05_created_at_brin_idx
ON events_2026_05
USING BRIN(created_at);
```

Run:

```bash
psql postgres://app:app@localhost:5432/system_design_postgres -f 08_partitioning_brin.sql
```
