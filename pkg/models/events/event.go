package events

import "time"

type Issuer string

const (
	IssuerAccount Issuer = "account"
	IssuerRobot   Issuer = "robot"
)

type Resource string

const (
	ResourceAccount   Resource = "account"
	ResourceContainer Resource = "container"
	ResourceImage     Resource = "image"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Metadata map[string]interface{}

func NewMetadata() *Metadata {
	return &Metadata{}
}

func (m *Metadata) WithReason(reason string) *Metadata {
	(*m)["reason"] = reason
	return m
}

func (m *Metadata) WithDescription(desc string) *Metadata {
	(*m)["description"] = desc
	return m
}

type Event struct {
	ID         string    `json:"id"`
	Resource   Resource  `json:"resource"`
	ResourceID string    `json:"resource_id"`
	Action     Action    `json:"action"`
	Issuer     Issuer    `json:"issuer"`
	IssuerID   string    `json:"issuer_id"`
	Metadata   Metadata  `json:"metadata"`
	CreatedAt  time.Time `json:"created_at"`
}
