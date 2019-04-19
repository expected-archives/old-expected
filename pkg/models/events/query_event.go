package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"time"
)

func eventFromRows(rows *sql.Rows) (*Event, error) {
	var metadata Metadata
	var eventIssuer, eventAction, eventResource, eventMetadata string

	event := &Event{}
	err := rows.Scan(&event.ID, &eventResource, &event.ResourceID, &eventAction, &eventIssuer, &event.IssuerID,
		&eventMetadata, &event.CreatedAt)
	if err != nil {
		return nil, err
	}
	event.Issuer = Issuer(eventIssuer)
	event.Action = Action(eventAction)
	event.Resource = Resource(eventResource)
	err = json.Unmarshal([]byte(eventMetadata), &metadata)
	if err != nil {
		return nil, err
	}
	event.Metadata = metadata
	return event, nil
}

func CreateEvent(ctx context.Context,
	res Resource, resourceId string, action Action, issuer Issuer, issuerId string, meta Metadata) error {
	event := Event{
		ID:         uuid.New().String(),
		Resource:   res,
		ResourceID: resourceId,
		Action:     action,
		Issuer:     issuer,
		IssuerID:   issuerId,
		Metadata:   meta,
		CreatedAt:  time.Now(),
	}
	strMetadata, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	_, err = services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO events VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, event.ID, event.Resource, event.ResourceID, event.Action, event.Issuer, event.IssuerID, strMetadata, event.CreatedAt)
	return err
}

func FindEventsByResourceID(ctx context.Context, resourceId string) ([]*Event, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
			SELECT id, resource, resource_id, action, issuer, issuer_id, metadata, created_at
			FROM events WHERE resource_id = $1
		`, resourceId)
	if err != nil {
		return nil, err
	}
	var events []*Event
	for rows.Next() {
		event, err := eventFromRows(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func FindEventsByIssuerID(ctx context.Context, issuerId string) ([]*Event, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
			SELECT id, resource, resource_id, action, issuer, issuer_id, metadata, created_at
			FROM events WHERE issuer_id = $1
		`, issuerId)
	if err != nil {
		return nil, err
	}
	var events []*Event
	for rows.Next() {
		event, err := eventFromRows(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
