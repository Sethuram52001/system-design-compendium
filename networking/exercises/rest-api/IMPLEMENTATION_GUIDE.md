# REST API Exercise: Implementation Guide

## Goal

Build a small Go HTTP API with Gin that demonstrates resource-oriented CRUD, JSON request/response handling, and HTTP status codes.

## Step 1: Create the Local Shape

Use this structure:

```text
networking/exercises/rest-api/
├── README.md
├── IMPLEMENTATION_GUIDE.md
└── solution/
    └── backend/
        ├── main.go
        ├── go.mod
        └── go.sum
```

## Step 2: Define the User Model

Use a small `User` struct:

- `id`
- `name`
- `email`

Keep the first version in memory with a map keyed by user ID.

## Step 3: Add Collection Routes

Implement:

- `GET /users`
- `POST /users`

Return `201 Created` for successful creates and `409 Conflict` if a caller tries to reuse an ID.

## Step 4: Add Resource Routes

Implement:

- `GET /users/:id`
- `PUT /users/:id`
- `PATCH /users/:id`
- `DELETE /users/:id`

Use `404 Not Found` when the user does not exist and `204 No Content` for successful deletes.

## Step 5: Test with curl

Run the server from `solution/backend`:

```bash
go run .
```

Then exercise the API:

```bash
curl http://localhost:8080/users
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

## Done Criteria

- All CRUD routes exist.
- JSON input and output work.
- Common error cases return clear JSON responses.
- Status codes match the behavior.
- Data is stored in memory for the first version.

