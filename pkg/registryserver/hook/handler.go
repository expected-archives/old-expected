package hook

import (
	"context"
	"github.com/docker/distribution/notifications"
)

// Handle trigger hook functions.
func Handle(envelope notifications.Envelope) error {
	for _, event := range envelope.Events {
		if event.Action == notifications.EventActionPush {
			if err := onPush(context.Background(), event); err != nil {
				return err
			}
		}
		if event.Action == notifications.EventActionDelete {
			if err := onDelete(context.Background(), event); err != nil {
				return err
			}
		}
	}
	return nil
}
