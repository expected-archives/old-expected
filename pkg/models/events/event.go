package events

import "time"

type Issuer string

const (
	IssuerAccount Issuer = "account"
	IssuerRobot   Issuer = "robot"
)

type Resource string

const (
	ResourceContainer  Resource = "container"
	ResourceImage      Resource = "image"
	ResourcePlan       Resource = "plan"
	ResourceCustomPlan Resource = "customplan"
	ResourceAccount    Resource = "account"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionDelete Action = "delete"
	ActionUpdate Action = "update"
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
