package response

import (
	"fmt"
	"net/http"
)

// Stream allows streaming a response back to the client using a callback.
func Stream(w http.ResponseWriter, callback func(w http.ResponseWriter)) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Transfer-Encoding", "chunked")
	
	// Create a wrapper to auto-flush on writes if needed, or rely on callback
	type flushWriter struct {
		http.ResponseWriter
		flusher http.Flusher
	}
	
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
