# HTTP and REST

## What It Is

HTTP is the application protocol most web systems use for request and response communication. A client sends a request to a server, and the server returns a response with a status code, headers, and usually a body.

REST is an API style built on top of HTTP. It models data as resources and uses HTTP methods to operate on those resources.

## Why It Exists

HTTP is simple, universal, and supported almost everywhere: browsers, mobile apps, backend services, proxies, gateways, caches, CDNs, and observability tools.

REST gives teams a common way to expose CRUD-like APIs without inventing a custom protocol for every service.

## Mental Model

Think in resources:

```text
/users
/users/123
/orders
/orders/abc
```

Then use HTTP methods to express intent:

| Method | Common Meaning |
|---|---|
| GET | Read |
| POST | Create or trigger an action |
| PUT | Replace |
| PATCH | Partially update |
| DELETE | Remove |

The server communicates the outcome with status codes such as `200`, `201`, `204`, `400`, `404`, and `500`.

## How It Works

A REST API usually combines:

- URL paths for resources
- HTTP methods for operations
- Headers for metadata
- JSON for request and response bodies
- Status codes for outcomes

Example:

```text
GET /users/123
```

returns one user, while:

```text
POST /users
```

creates one.

## When To Use It

Use HTTP/REST when:

- you are building public APIs
- browser/mobile/client compatibility matters
- resources map naturally to CRUD operations
- you want easy testing with curl, Postman, or browser tools
- operational simplicity matters more than maximum efficiency

## When Not To Use It

REST may be a weaker fit when:

- clients need flexible nested queries
- internal services need strict generated contracts
- low-latency binary communication matters
- the workflow is naturally streaming or bidirectional

GraphQL, gRPC, SSE, or WebSockets may fit those cases better.

## Common Pitfalls

- Treating every action as `POST /doThing` and losing resource clarity.
- Returning `200 OK` for every outcome instead of meaningful status codes.
- Designing endpoints around database tables instead of user-facing resources.
- Ignoring idempotency for retries.
- Returning inconsistent error shapes.

## Related Concepts

- GraphQL
- gRPC
- Server-Sent Events
- Caching
- Load balancing

## Exercise

- [REST API Exercise](../exercises/rest-api/README.md)

## Resources

- [Networking resources](../resources.md)

