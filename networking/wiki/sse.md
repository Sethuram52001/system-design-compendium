# Server-Sent Events

## What It Is

Server-Sent Events, or SSE, is a browser-supported pattern for streaming events from a server to a client over HTTP.

The client opens a connection, the server keeps the response open, and the server writes events over time using the `text/event-stream` format.

## Why It Exists

Many realtime features are one-way: the server needs to push updates to the client, but the client does not need to constantly send messages back on the same connection.

SSE keeps this case simple. It works over normal HTTP and has a built-in browser API called `EventSource`.

## Mental Model

SSE is a long-lived HTTP response:

```text
client -> GET /chat/stream
server -> event: token
server -> event: token
server -> event: done
```

The response does not complete immediately. The server writes chunks and flushes them as events become available.

## How It Works

The server responds with:

```http
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

Then it writes events:

```text
event: token
data: {"delta":"hello"}

```

The blank line ends the event.

Common SSE fields:

| Field | Purpose |
|---|---|
| `event:` | Names the event type |
| `data:` | Carries the event payload |
| `id:` | Sets the event ID for reconnects |
| `retry:` | Suggests reconnect delay in milliseconds |
| `:` | Comment line, often used as a heartbeat |

In browsers, `EventSource` parses these fields. Named events are delivered to matching event listeners. `retry`, `id`, and comments are handled by the browser internally.

## When To Use It

Use SSE when:

- updates are server-to-client
- the client sends one request and receives a stream
- simple HTTP infrastructure is preferred
- automatic browser reconnect behavior is useful
- examples include LLM response streaming, progress updates, notifications, live logs, dashboards, and activity feeds

## When Not To Use It

Avoid SSE when:

- communication must be bidirectional
- the client needs to send frequent messages over the same connection
- binary messages are central
- peer-to-peer communication is needed
- WebSocket or WebRTC semantics fit better

## Common Pitfalls

- Forgetting to flush after writing each event.
- Missing the blank line after `data:`.
- Assuming heartbeat comments trigger JavaScript handlers.
- Closing the `EventSource` inside `onerror`, which disables automatic reconnect.
- Forgetting that proxies may buffer or close idle connections.
- Relying on custom headers with native browser `EventSource`; cookies or other auth strategies are usually simpler.

## Related Concepts

- HTTP/REST
- WebSockets
- WebRTC
- Streaming APIs
- Load balancer idle timeouts

## Exercise

- [Server-Sent Events Exercise](../exercises/sse/README.md)

## Resources

- [Networking resources](../resources.md)

