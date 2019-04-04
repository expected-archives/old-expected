package hook

import (
	"context"
	"github.com/docker/distribution/notifications"
)

// Handle trigger hook functions.
func Handle(envelope notifications.Envelope) error {
	for _, event := range envelope.Events {
		if event.Action == "push" {
			if err := Push(context.Background(), event); err != nil {
				return err
			}
		}
	}
	return nil
}
