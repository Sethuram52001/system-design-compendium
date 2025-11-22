# gRPC Exercise: Inter-Service Communication

## Exercise Description

In this exercise, you will build two separate Go servers communicating with each other using gRPC. This simulates a common backend scenario where microservices interact efficiently via RPC.

- Implement a **UserService** gRPC server that manages user data in-memory.
- Implement a **ClientService** gRPC client server that calls the UserService to fetch user data.
- Explore defining gRPC service, messages via Protocol Buffers (.proto files).
- Generate Go code from proto files and use the Go gRPC library to implement servers and clients.
- Focus on unary RPC calls for simplicity.

## Requirements

- Define a gRPC service `UserService` with a unary RPC method:  
  - `GetUser(GetUserRequest) returns (GetUserResponse)`  
- `UserService` server stores user info in-memory (map keyed by user ID).  
- `ClientService` server acts as a gRPC client to call `UserService.GetUser` and expose an HTTP endpoint (`/user/{id}`) to fetch user data via REST forwarding.  
- Implement basic error handling and proper gRPC status codes.  
- Use Go modules and standard gRPC packages.

## Project Structure Suggestion

exercises/grpc/
├── user-service/             # gRPC server managing user data
│   ├── server.go
│   ├── userpb/               # Generated protobuf Go files
│   └── users.go              # User management and storage logic
│
├── client-service/           # gRPC client acting as API gateway
│   ├── server.go             # HTTP server forwarding to gRPC client calls
│   ├── userpb/               # Same protobuf generated files
│   └── client.go             # gRPC client implementation
│
├── proto/
│   └── user.proto            # Service and message definitions
│
└── README.md                 # This exercise documentation


## Learning Objectives

- Understand and implement gRPC services and messages with Protocol Buffers.  
- Learn how to generate Go code from proto definitions and use it in server and client.  
- Explore service-to-service communication using unary RPC calls.  
- Combine gRPC client calls with HTTP server endpoints for easy testing.  
- Gain knowledge of error handling and status codes in gRPC.

## How to Run

1. Generate Go gRPC code from the `.proto` file using protoc with Go plugins.  
2. Start the UserService gRPC server.  
3. Start the ClientService HTTP server which acts as a gRPC client.  
4. Use curl or HTTP clients to call ClientService’s `/user/{id}` endpoint to fetch data, which internally calls UserService over gRPC.
