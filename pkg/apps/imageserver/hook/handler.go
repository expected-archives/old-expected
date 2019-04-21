package hook

import (
	"context"
	"errors"
	"github.com/docker/distribution/notifications"
	"strings"
)

// Handle trigger hook functions.
func Handle(envelope notifications.Envelope) error {
	for _, event := range envelope.Events {
		if event.Action == notifications.EventActionPush {
			if err := onPush(context.Background(), event); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseRepository return the namespace id and the name of the image.
// Can throw an error only if the repository is malformed.
func parseRepository(repo string) (namespaceID, name string, err error) {
	str := strings.Split(repo, "/")
	if len(str) != 2 {
		return "", "", errors.New("repository is malformed")
	}
	return str[0], str[1], nil
}
