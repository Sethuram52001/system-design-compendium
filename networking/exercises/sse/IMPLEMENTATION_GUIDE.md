# Server-Sent Events Exercise: Implementation Guide

## Goal

Build the smallest useful SSE demo for a modern LLM-style use case: a Go server exposes `GET /chat/stream`, keeps the HTTP connection open, and streams a fake assistant response one chunk at a time.

No AI endpoint is called. The point is to understand the streaming transport.

## Mental Model

SSE is just HTTP where the response does not finish immediately.

Instead of:

```text
client -> request
server -> one response
connection closes
```

SSE does this:

```text
client -> request /chat/stream?message=...
server -> token event
server -> token event
server -> token event
server -> done event
connection stays open
```

This is a natural fit for LLM-style responses because the client usually sends one prompt, then the server streams output back in one direction.

## SSE Wire Format

The simplest valid event looks like this:

```text
data: hello

```

The blank line at the end matters. It tells the client that one event is complete.

Useful fields:

- `data:` event payload
- `event:` optional event name
- `id:` optional event ID for reconnects
- `retry:` optional reconnect delay in milliseconds

Example:

```text
id: 1
event: token
data: {"delta":"Hello"}

```

## Step 1: Create the Solution Folder

From `networking/exercises/sse`:

```bash
mkdir -p solution/server
cd solution/server
go mod init github.com/Sethuram52001/system-design-compendium/networking/exercises/sse/solution/server
```

## Step 2: Create the Simulated LLM Stream Handler

Create `main.go`.

The handler needs to:

- Confirm the response supports flushing.
- Set SSE headers.
- Read the user's input message from the query string.
- Choose a fake assistant response.
- Split the fake response into small chunks.
- Write each chunk as a `token` event.
- Write a final `done` event.
- Flush after each event.
- Stop when the client disconnects.

Use this structure:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type tokenEvent struct {
	Delta string `json:"delta"`
}

type doneEvent struct {
	Reason string `json:"reason"`
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "event: %s\n", event); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
		return err
	}

	flusher.Flush()
	return nil
}

func fakeLLMResponse(message string) []string {
	response := "You asked: " + message + ". In a real LLM API, the model would generate this response incrementally. Here, we split a canned response into chunks so the client can learn how SSE streaming works."
	return strings.Split(response, " ")
}

func chatStreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	message := r.URL.Query().Get("message")
	if message == "" {
		message = "Explain Server-Sent Events"
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, chunk := range fakeLLMResponse(message) {
		select {
		case <-r.Context().Done():
			log.Println("client disconnected")
			return
		default:
			if err := writeSSE(w, flusher, "token", tokenEvent{Delta: chunk + " "}); err != nil {
				log.Printf("failed to write token: %v", err)
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
	}

	if err := writeSSE(w, flusher, "done", doneEvent{Reason: "stop"}); err != nil {
		log.Printf("failed to write done event: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!doctype html>
<html>
<body>
  <h1>Simulated LLM SSE Stream</h1>
  <form id="form">
    <input id="message" value="Explain SSE in one sentence" />
    <button type="submit">Stream</button>
  </form>
  <pre id="output"></pre>
  <script>
    const form = document.getElementById("form");
    const message = document.getElementById("message");
    const output = document.getElementById("output");

    form.addEventListener("submit", (event) => {
      event.preventDefault();
      output.textContent = "";

      const stream = new EventSource("/chat/stream?message=" + encodeURIComponent(message.value));

      stream.addEventListener("token", (event) => {
        const payload = JSON.parse(event.data);
        output.textContent += payload.delta;
      });

      stream.addEventListener("done", () => {
        stream.close();
      });

      stream.onerror = () => {
        stream.close();
      };
    });
  </script>
</body>
</html>`)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/chat/stream", chatStreamHandler)

	log.Println("SSE server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Step 3: Run It

```bash
go run .
```

Open:

```text
http://localhost:8080
```

Type a message and submit the form. You should see the fake assistant response appear chunk by chunk.

You can also test with curl:

```bash
curl -N "http://localhost:8080/chat/stream?message=Explain%20SSE"
```

The `-N` flag disables curl buffering so streamed events show up as they arrive.

## Step 4: Explain the Important Pieces

- `text/event-stream` tells the client this is an SSE stream.
- `http.Flusher` lets the server push buffered data immediately.
- `r.Context().Done()` tells the server when the client disconnected.
- `EventSource` is the browser API for consuming SSE.
- `event: token` lets the client handle streamed chunks separately from completion metadata.
- `event: done` tells the client to close the stream.
- SSE automatically reconnects if the connection drops.

## Done Criteria

- `GET /chat/stream?message=...` streams fake LLM response chunks.
- Browser client receives `token` and `done` events with `EventSource`.
- `curl -N` shows streamed event output.
- Server logs client disconnects.
- You can explain when SSE is enough and when WebSockets are needed.

## When To Use SSE

Use SSE when:

- updates flow from server to client
- the client does not need to send frequent messages back over the same connection
- you want simple browser support over HTTP
- examples include LLM response streaming, notifications, progress updates, dashboards, logs, and feeds

## When Not To Use SSE

Avoid SSE when:

- you need bidirectional realtime communication
- clients need to send frequent messages to the server
- you need peer-to-peer browser communication
- binary messages are central to the protocol

Use WebSockets or WebRTC for those cases.
