package containers

import (
	"time"
)

type State string

const (
	StateStopped  State = "stopped"
	StateStarting State = "starting"
	StateRunning  State = "running"
)

type Container struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Endpoints   []*Endpoint       `json:"endpoint"`
	PlanID      string            `json:"plan_id"`
	Environment map[string]string `json:"environment"`
	Tags        []string          `json:"tags"`
	NamespaceID string            `json:"namespace_id"`
	State       State             `json:"state"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type Endpoint struct {
	ID        string    `json:"id"`
	Endpoint  string    `json:"endpoint"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"created_at"`
}
