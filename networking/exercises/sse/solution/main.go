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
	response := "You asked: " + message + ". In a real LLM API , the model would generate the response incrementally, allowing us to stream tokens as they are produced.. Here, we split a canned responsed into chunks so the client can learn how SSE streaming works."

	return strings.Split(response, " ")
}

func chatStreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
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
				log.Printf("failed to write token event: %v", err)
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
