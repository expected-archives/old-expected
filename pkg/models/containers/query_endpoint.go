package containers

import (
	"context"
	"database/sql"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"time"
)

func endpointFromRows(rows *sql.Rows) (*Endpoint, error) {
	endpoint := &Endpoint{}
	err := rows.Scan(&endpoint.ID, &endpoint.Endpoint, &endpoint.Default, &endpoint.CreatedAt)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func CreateEndpoint(ctx context.Context, container *Container, endpoint string, isDefault bool) (*Endpoint, error) {
	id := uuid.New().String()
	createdAt := time.Now()
	if _, err := services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO containers_endpoints (id, container_id, endpoint, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, container.ID, endpoint, isDefault, createdAt); err != nil {
		return nil, err
	}
	e := &Endpoint{
		ID:        id,
		Endpoint:  endpoint,
		Default:   isDefault,
		CreatedAt: createdAt,
	}
	container.Endpoints = append(container.Endpoints, e)
	return e, nil
}

func UpdateEndpoint(ctx context.Context, endpoint *Endpoint) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		UPDATE containers_endpoints SET endpoint = $2, is_default = $3, created_at = $4
		WHERE id = $1
	`, endpoint.ID, endpoint.Endpoint, endpoint.Default, endpoint.CreatedAt)
	return err
}

func DeleteEndpoint(ctx context.Context, id string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		DELETE FROM containers_endpoints WHERE id = $1
	`, id)
	return err
}

func FindEndpointsByContainerID(ctx context.Context, id string) ([]*Endpoint, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, endpoint, is_default, created_at
		FROM containers_endpoints WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var endpoints []*Endpoint
	for rows.Next() {
		endpoint, err := endpointFromRows(rows)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}
