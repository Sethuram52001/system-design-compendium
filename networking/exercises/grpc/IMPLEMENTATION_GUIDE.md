# gRPC Exercise: Implementation Guide

## Goal

Build two small Go services: a `UserService` gRPC server and a `ClientService` HTTP gateway that calls it.

## Step 1: Define the Contract

Start with `proto/user.proto` and define:

- `User`
- `GetUserRequest`
- `GetUserResponse`
- `UserService.GetUser`

The `.proto` file is the shared contract between services.

## Step 2: Generate Go Code

From `networking/exercises/grpc`, generate protobuf and gRPC bindings:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto
```

Copy generated files into each solution service package if you keep each service as a separate Go module.

## Step 3: Implement UserService

Inside `solution/user-service`:

- Create an in-memory user store.
- Implement `GetUser`.
- Return `InvalidArgument` when the ID is empty.
- Return `NotFound` when the user does not exist.
- Listen on port `50051`.

## Step 4: Implement ClientService

Inside `solution/client-service`:

- Connect to `localhost:50051`.
- Expose `GET /user/:id` over HTTP.
- Call `UserService.GetUser`.
- Convert gRPC status codes into HTTP status codes.

## Step 5: Run Locally

Start the gRPC service first, then the HTTP gateway:

```bash
cd solution/user-service
go run .
```

```bash
cd solution/client-service
go run .
```

Then test:

```bash
curl http://localhost:8080/user/1
```

## Done Criteria

- The proto contract is clear.
- UserService returns correct gRPC status codes.
- ClientService forwards HTTP requests to gRPC.
- Missing users return HTTP `404`.
- Empty or invalid IDs return HTTP `400`.

