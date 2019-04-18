package sse

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPHandler serves new connections with events for a given stream ...
func (s *Server) HTTPHandler(streamID string, w http.ResponseWriter, r *http.Request) {
	flusher, err := w.(http.Flusher)
	if !err {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stream := s.getStream(streamID)

	if stream == nil && !s.AutoStream {
		http.Error(w, "Stream not found!", http.StatusInternalServerError)
		return
	} else if stream == nil && s.AutoStream {
		stream = s.CreateStream(streamID)
	}

	eventid := r.Header.Get("Last-Event-ID")
	if eventid == "" {
		eventid = "0"
	}

	// Create the stream subscriber
	sub := stream.addSubscriber(eventid)
	defer sub.close()

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		sub.close()
	}()

	// Push events to client
	for {
		select {
		case ev, ok := <-sub.connection:
			if !ok {
				return
			}

			// If the data buffer is an empty string abort.
			if len(ev.Data) == 0 {
				break
			}

			// if the event has expired, dont send it
			if s.EventTTL != 0 && time.Now().After(ev.timestamp.Add(s.EventTTL)) {
				continue
			}

			fmt.Fprintf(w, "id: %s\n", ev.ID)
			fmt.Fprintf(w, "data: %s\n", ev.Data)
			if len(ev.Event) > 0 {
				fmt.Fprintf(w, "event: %s\n", ev.Event)
			}
			if len(ev.Retry) > 0 {
				fmt.Fprintf(w, "retry: %s\n", ev.Retry)
			}
			fmt.Fprint(w, "\n")
			flusher.Flush()
		}
	}
}