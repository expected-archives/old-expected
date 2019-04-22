package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SetupConnection(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

func SendJSON(w http.ResponseWriter, event string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Send(w, event, string(b))
}

func Send(w http.ResponseWriter, event, data string) error {
	flusher, _ := w.(http.Flusher)
	if _, err := fmt.Fprintf(w, "event: %s\n", event); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "data: %s\n", data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}
