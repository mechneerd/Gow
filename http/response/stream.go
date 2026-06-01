package response

import (
	"fmt"
	"net/http"
	"time"
)

// flushWriter wraps an http.ResponseWriter to auto-flush.
type flushWriter struct {
	http.ResponseWriter
	flusher http.Flusher
}

// Stream allows streaming a response back to the client using a callback.
func Stream(w http.ResponseWriter, callback func(w http.ResponseWriter)) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Transfer-Encoding", "chunked")
	
	fw := &flushWriter{w, flusher}

	callback(fw)
	return nil
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.ResponseWriter.Write(p)
	if err == nil {
		fw.flusher.Flush()
	}
	return n, err
}

// Download sends a file as a stream for downloading.
func Download(w http.ResponseWriter, r *http.Request, filePath string, fileName string) {
	if fileName != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	}
	http.ServeFile(w, r, filePath)
}

// StreamLines streams content line by line with a callback.
func StreamLines(w http.ResponseWriter, callback func(w http.ResponseWriter) error) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fw := &flushWriter{w, flusher}
	return callback(fw)
}

// SSEEvent represents a Server-Sent Event.
type SSEEvent struct {
	ID    string
	Event string
	Data  string
	Retry int
}

// StreamSSE streams Server-Sent Events to the client.
func StreamSSE(w http.ResponseWriter, events <-chan SSEEvent) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	for event := range events {
		if event.ID != "" {
			fmt.Fprintf(w, "id: %s\n", event.ID)
		}
		if event.Event != "" {
			fmt.Fprintf(w, "event: %s\n", event.Event)
		}
		if event.Retry > 0 {
			fmt.Fprintf(w, "retry: %d\n", event.Retry)
		}
		fmt.Fprintf(w, "data: %s\n\n", event.Data)
		flusher.Flush()
	}

	return nil
}

// StreamTimer streams events at regular intervals.
func StreamTimer(w http.ResponseWriter, interval time.Duration, callback func(int) string) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	i := 0
	for range ticker.C {
		data := callback(i)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		i++
	}

	return nil
}

