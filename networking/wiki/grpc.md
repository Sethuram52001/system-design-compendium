# gRPC

## What It Is

gRPC is an RPC framework where services expose typed remote methods. Instead of designing URLs and JSON payloads by hand, you define a service contract in a `.proto` file and generate client/server code from it.

gRPC commonly uses Protocol Buffers for serialization and HTTP/2 as the transport.

## Why It Exists

REST is easy and universal, but internal service-to-service communication often benefits from stronger contracts, generated clients, smaller payloads, and efficient transport.

gRPC exists to make remote calls feel like typed method calls while keeping the contract explicit and language-neutral.

## Mental Model

The `.proto` file is the API contract:

```proto
service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

From that contract, tooling generates:

- request and response structs
- a client interface
- a server interface
- registration and wiring helpers

Your application implements the generated server interface. Other services call it through the generated client.

## How It Works

Typical flow:

```text
user.proto
  -> protoc generates Go code
  -> server implements generated interface
  -> client calls generated client method
  -> gRPC handles network transport and serialization
```

For a unary RPC, the shape is similar to a normal function call:

```text
GetUser(request) -> response
```

But the call crosses the network.

## When To Use It

Use gRPC when:

- services communicate internally
- strict contracts matter
- many languages need generated clients
- payload efficiency matters
- low-latency service-to-service calls are important
- streaming RPCs may be useful later

## When Not To Use It

gRPC may be a weaker fit when:

- the API is public and browser-first
- easy manual testing is important
- clients cannot easily use generated code
- simple HTTP/JSON is good enough
- infrastructure does not support HTTP/2 well

## Common Pitfalls

- Forgetting that generated files should not be manually edited.
- Treating protobuf field numbers as disposable.
- Hiding business errors inside generic internal errors.
- Not setting timeouts or deadlines on client calls.
- Making every service call synchronous without thinking about failure chains.

## Related Concepts

- HTTP/REST
- Protocol Buffers
- HTTP/2
- Service discovery
- Retries and timeouts

## Exercise

- [gRPC Exercise](../exercises/grpc/README.md)

## Resources

- [Networking resources](../resources.md)

