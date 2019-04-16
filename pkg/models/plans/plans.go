package plans

import "time"

type Type string

const (
	Container Type = "container"
	Image          = "registry"
)

type Metadata map[string]interface{}

type Plan struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Price     float32   `json:"price"`
	Metadata  Metadata  `json:"metadata"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
