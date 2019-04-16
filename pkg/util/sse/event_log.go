package sse

import (
	"strconv"
	"time"
)

// EventLog holds all of previous events
type EventLog []*Event

// Add event to eventlog
func (e *EventLog) Add(ev *Event) {
	ev.ID = []byte(e.currentindex())
	ev.timestamp = time.Now()
	*e = append(*e, ev)
}

// Clear events from eventlog
func (e *EventLog) Clear() {
	*e = nil
}

// Replay events to a subscriber
func (e *EventLog) Replay(s *Subscriber) {
	for i := 0; i < len(*e); i++ {
		if string((*e)[i].ID) >= s.eventid {
			s.connection <- (*e)[i]
		}
	}
}

func (e *EventLog) currentindex() string {
	return strconv.Itoa(len(*e))
}
