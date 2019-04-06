package hook

import (
	"context"
	"github.com/docker/distribution/notifications"
)

func onDelete(ctx context.Context, event notifications.Event) error {
	return nil
}
