package sse

// Subscriber ...
type Subscriber struct {
	eventid    string
	quit       chan *Subscriber
	connection chan *Event
}

// Close will let the stream know that the clients connection has terminated
func (s *Subscriber) close() {
	s.quit <- s
}
