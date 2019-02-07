package models

import "time"

type ContainersModel struct{}

type Container struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Image     string    `json:"image"`
	Endpoint  string    `json:"endpoint"`
	Memory    int       `json:"memory"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
