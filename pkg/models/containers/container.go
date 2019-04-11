package containers

import (
	"time"
)

type Container struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Endpoint    string            `json:"endpoint"`
	Memory      int               `json:"memory"`
	Environment map[string]string `json:"environment"`
	Tags        []string          `json:"tags"`
	OwnerID     string            `json:"-"`
	CreatedAt   time.Time         `json:"created_at"`
}
