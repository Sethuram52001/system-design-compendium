# Networking

## What This Topic Covers

Networking is the foundation for distributed systems: independent machines need to discover each other, establish connections, exchange data, handle failures, and route traffic as the system grows. This topic starts with the parts that show up most often in system design: network layers, DNS, TCP/UDP, HTTP, REST, GraphQL, gRPC, realtime protocols, and load balancing.

The goal is not to memorize every detail of networking. The goal is to build enough intuition to choose the right communication pattern, reason about latency, and explain how requests move through a real system.

## Why It Matters

Every system design eventually becomes a networking problem. Once a service has clients, dependencies, multiple replicas, multiple regions, or realtime features, design choices depend on network behavior: round trips, connection setup, packet loss, retries, timeouts, persistent connections, routing, and load balancing.

For interviews and practical engineering, networking knowledge helps answer questions like:

- Should this be REST, GraphQL, gRPC, SSE, WebSockets, or WebRTC?
- Is reliability more important than latency for this path?
- What happens when a backend server goes down?
- Where do load balancers fit?
- What state is held by clients, servers, or connections?
- How does this design behave across regions?

## Core Ideas

- **Layering:** IP routes packets, TCP/UDP/QUIC provide transport behavior, and application protocols such as HTTP, DNS, WebSockets, and gRPC define how applications communicate.
- **Discovery:** DNS turns names into addresses and can also participate in coarse load distribution through returned IPs and TTLs.
- **Connection setup:** Protocols have setup costs. TCP uses a handshake; HTTPS also adds TLS negotiation; persistent connections and multiplexing reduce repeated setup overhead.
- **Reliability vs latency:** TCP provides ordered, reliable delivery with overhead. UDP is lower overhead but leaves reliability and ordering to the application.
- **Request/response vs streaming:** HTTP APIs are often request/response; SSE, WebSockets, gRPC streaming, and WebRTC support longer-lived or realtime communication patterns.
- **Public vs internal APIs:** Public APIs usually favor simple, interoperable protocols such as HTTP/JSON. Internal service-to-service calls can use more efficient contracts like gRPC.
- **Load balancing:** Client-side, DNS-based, Layer 4, and Layer 7 load balancing all route traffic differently and have different tradeoffs.

## Mental Model

Think of networking as a stack of promises. Each layer gives the layer above it a simpler abstraction:

- IP gives you addressing and best-effort routing.
- TCP gives you a reliable ordered stream.
- UDP gives you fast datagrams with fewer guarantees.
- TLS gives encrypted transport.
- HTTP gives request/response semantics.
- REST, GraphQL, and gRPC give application-level contracts.

The higher you go in the stack, the easier the programming model becomes, but the more overhead and policy you usually add. System design is about choosing the simplest protocol that satisfies the product and scale requirements.

## Common Tradeoffs

| Choice | Good For | Tradeoff |
|---|---|---|
| HTTP/REST | Public APIs, CRUD resources, broad client support | Can over-fetch or under-fetch; JSON adds serialization overhead |
| GraphQL | Flexible clients, complex frontend data needs, changing product requirements | Resolver complexity, backend latency risks, harder caching |
| gRPC | Internal service-to-service calls, strong contracts, efficient serialization | Less friendly for public/browser clients |
| TCP | Reliable ordered delivery | Connection setup and retransmission overhead |
| UDP | Low-latency media, games, telemetry, DNS-style lookups | No built-in delivery or ordering guarantees |
| SSE | Server-to-client event streams over HTTP | One-way push; long-lived connections need care |
| WebSockets | Bidirectional realtime communication | Stateful connections complicate scaling and load balancing |
| WebRTC | Peer-to-peer media and browser communication | NAT traversal and fallback complexity |
| Layer 4 load balancing | Persistent connections, raw performance | Cannot route based on HTTP paths or headers |
| Layer 7 load balancing | HTTP routing, path/header/cookie-based routing | More processing overhead |

## Learning Objectives

- Understand common service communication styles.
- Compare REST, GraphQL, and gRPC tradeoffs.
- Practice explicit contracts between clients and servers.
- Build intuition for payload shape, status codes, and error handling.
- Create a foundation for later topics such as load balancing, retries, timeouts, and service discovery.

## Exercises

- [REST API](exercises/rest-api/README.md): HTTP methods, status codes, JSON payloads, and resource-oriented CRUD.
- [GraphQL](exercises/graphql/README.md): schema-first APIs, queries, mutations, and resolver behavior.
- [gRPC](exercises/grpc/README.md): protobuf contracts, unary RPC, status codes, and service-to-service communication.
- [Server-Sent Events](exercises/sse/README.md): one-way server-to-client realtime updates over HTTP.

## Wiki

- [HTTP and REST](wiki/rest.md)
- [gRPC](wiki/grpc.md)
- [Server-Sent Events](wiki/sse.md)

## Resources

- [Networking resources](resources.md)

## Related Topics

- Load balancing
- Caching
- Databases
- Observability
- Consistency
