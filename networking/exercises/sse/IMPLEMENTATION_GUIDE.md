# Server-Sent Events Exercise: Implementation Guide

## Goal

Build a small SSE endpoint that simulates an LLM response streaming back one chunk at a time.

No AI endpoint is called. The exercise is about the transport pattern:

```text
client -> GET /chat/stream?message=...
server -> token event
server -> token event
server -> done event
```

## Mental Model

SSE is a long-lived HTTP response. The server writes text in a specific format and flushes each event so the client receives it immediately.

The most important format rule is the blank line:

```text
event: token
data: {"delta":"hello"}

```

That blank line marks the end of one event.

## Step 1: Create the Solution

From `networking/exercises/sse`:

```bash
mkdir -p solution
cd solution
go mod init github.com/Sethuram52001/system-design-compendium/networking/exercises/sse/solution/server
```

Create `main.go`.

## Step 2: Add the Event Types

Use small structs for the event payloads:

```go
type tokenEvent struct {
	Delta string `json:"delta"`
}

type doneEvent struct {
	Reason string `json:"reason"`
}
```

`tokenEvent` carries one response chunk. `doneEvent` tells the client the stream is finished.

## Step 3: Add the SSE Writer

This helper writes one SSE event and flushes it:

```go
func writeSSE(w http.ResponseWriter, flusher http.Flusher, id int, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "id: %d\n", id)
	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", payload)
	flusher.Flush()
	return nil
}
```

Important pieces:

- `id:` gives the event an ID for reconnect behavior.
- `event:` names the event type, such as `token` or `done`.
- `data:` carries the payload.
- `\n\n` ends the event.
- `Flush()` pushes buffered bytes to the client immediately.

## Step 4: Add Retry and Heartbeat Helpers

`retry:` suggests how long the browser should wait before reconnecting:

```go
func writeRetry(w http.ResponseWriter, flusher http.Flusher, retryMillis int) {
	fmt.Fprintf(w, "retry: %d\n\n", retryMillis)
	flusher.Flush()
}
```

Heartbeat comments keep idle connections alive:

```go
func writeHeartbeat(w http.ResponseWriter, flusher http.Flusher) {
	fmt.Fprint(w, ": heartbeat\n\n")
	flusher.Flush()
}
```

Lines starting with `:` are SSE comments. The browser ignores them as events, but the bytes still keep the connection active.

## Step 5: Add the Stream Handler

The handler should:

- check that streaming is supported
- set SSE headers
- read `message` from the query string
- send `retry: 3000`
- stream fake tokens
- send periodic heartbeats
- send `done`
- stop if the client disconnects

Core shape:

```go
func chatStreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	message := r.URL.Query().Get("message")
	if message == "" {
		message = "Explain Server-Sent Events"
	}

	writeRetry(w, flusher, 3000)

	eventID := 1
	for index, chunk := range fakeLLMResponse(message) {
		select {
		case <-r.Context().Done():
			return
		default:
			if index > 0 && index%8 == 0 {
				writeHeartbeat(w, flusher)
			}

			writeSSE(w, flusher, eventID, "token", tokenEvent{Delta: chunk + " "})
			eventID++
			time.Sleep(250 * time.Millisecond)
		}
	}

	writeSSE(w, flusher, eventID, "done", doneEvent{Reason: "stop"})
}
```

Your fake response function can be simple:

```go
func fakeLLMResponse(message string) []string {
	response := "You asked: " + message + ". This canned response is split into chunks to simulate LLM streaming."
	return strings.Split(response, " ")
}
```

## Step 6: Wire the Server

For a curl-first exercise, you only need:

```go
func main() {
	http.HandleFunc("/chat/stream", chatStreamHandler)

	log.Println("SSE server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

A tiny browser UI with `EventSource` is optional. The backend concept can be tested fully with curl.

## Step 7: Test the Stream

Run:

```bash
go run .
```

Then:

```bash
curl -N "http://localhost:8080/chat/stream?message=Explain%20SSE"
```

Expected shape:

```text
retry: 3000

id: 1
event: token
data: {"delta":"You "}

: heartbeat

id: 2
event: token
data: {"delta":"asked: "}
```

The `-N` flag disables curl buffering so events print as they arrive.

## Step 8: Test Last-Event-ID

Browsers send `Last-Event-ID` after reconnecting. You can simulate it manually:

```bash
curl -N -H "Last-Event-ID: 5" "http://localhost:8080/chat/stream?message=Explain%20SSE"
```

This exercise only shows the header behavior. It does not implement true resume logic yet.

## Key Things To Remember

- `text/event-stream` tells the client this is SSE.
- `data:` is the only required SSE field.
- `event:` lets clients listen for named events.
- `id:` helps reconnect/resume behavior.
- `retry:` configures automatic browser reconnect delay.
- `: heartbeat` is ignored by the client but keeps the connection alive.
- `Flush()` is what makes streaming visible immediately.
- SSE is server-to-client only. Use WebSockets when the client must send frequent messages back over the same connection.

## Done Criteria

- `GET /chat/stream?message=...` streams fake LLM response chunks.
- The stream includes `retry:`, `id:`, `event:`, and `data:` lines.
- Heartbeat comments appear during the stream.
- `curl -N` shows events as they arrive.
- You can explain why SSE fits LLM response streaming.

