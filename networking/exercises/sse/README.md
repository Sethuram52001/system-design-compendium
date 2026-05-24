# Server-Sent Events Exercise: Simulated LLM Response Streaming

## Exercise Description

Build a tiny Server-Sent Events server that simulates an LLM response streaming token-by-token to a client. This exercise does not call any AI endpoint. It uses a hardcoded response and streams it in small chunks so you can understand the transport pattern used by many chat and completion APIs.

You will:

- Create an HTTP endpoint that keeps the connection open.
- Send events using the `text/event-stream` format.
- Stream a fake assistant response one chunk at a time.
- Include a final `done` event when the response is complete.
- Test the stream with curl or a tiny browser client.
- Understand where SSE fits compared with WebSockets and polling.

## Tech Stack

- Go
- Standard `net/http`
- curl
- Optional browser `EventSource`

## Requirements

- Implement `GET /chat/stream?message=...` as an SSE stream.
- Set the correct response headers:
  - `Content-Type: text/event-stream`
  - `Cache-Control: no-cache`
  - `Connection: keep-alive`
- Send fake response chunks using an event such as `token`.
- Send a final event such as `done`.
- Include event data using JSON strings.
- Test with `curl -N`.
- Handle client disconnects cleanly.

## Extension: Resume and Keepalive Behavior

Extend the basic stream with production-adjacent SSE features:

- Send `id:` on each `token` and `done` event.
- Send `retry: 3000` when the stream starts.
- Send heartbeat comments such as `: heartbeat` during the stream.
- Test the raw stream with `curl -N` and confirm those lines appear.
- Explain what the browser does with `Last-Event-ID` after reconnecting.

## Learning Objectives

- Understand SSE as one-way server push over HTTP.
- Learn the SSE wire format.
- Practice long-lived HTTP connections.
- Understand why SSE works well for streaming LLM-style responses.
- Learn how event IDs, retry hints, and heartbeat comments make long-lived streams more resilient.
- Learn where SSE is useful: chat completions, progress updates, notifications, dashboards, feeds, and logs.

## Reference Shape

```text
networking/exercises/sse/
├── README.md
├── IMPLEMENTATION_GUIDE.md
└── solution/
    └── server/
        ├── main.go
        └── go.mod
```

## Stretch Goals

- Add named events: `event: token`, `event: metadata`, `event: usage`, and `event: done`.
- Add a fake `usage` event with token counts.
- Stream different canned responses based on the input message.
- Compare SSE behavior with returning the full response at once.
