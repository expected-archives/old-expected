package containers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"time"
)

func containerFromRows(rows *sql.Rows) (*Container, error) {
	var environmentJson, tagsJson string
	container := &Container{}
	err := rows.Scan(&container.ID, &container.Name, &container.Image, &container.PlanID, &environmentJson, &tagsJson,
		&container.NamespaceID, &container.State, &container.CreatedAt, &container.UpdatedAt)
	if err != nil {
		return nil, err
	}
	var environment map[string]string
	if err = json.Unmarshal([]byte(environmentJson), &environment); err != nil {
		return nil, err
	}
	container.Environment = environment
	var tags []string
	if err = json.Unmarshal([]byte(tagsJson), &tags); err != nil {
		return nil, err
	}
	container.Tags = tags
	endpoints, err := FindEndpointsByContainerID(context.Background(), container.ID)
	if err != nil {
		return nil, err
	}
	if endpoints == nil {
		endpoints = []*Endpoint{}
	}
	container.Endpoints = endpoints
	return container, nil
}

func CreateContainer(ctx context.Context, name, image, planId string, environment map[string]string,
	tags []string, namespaceId string) (*Container, error) {
	jsonEnvironment, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}
	container := &Container{
		ID:          uuid.New().String(),
		Name:        name,
		Image:       image,
		PlanID:      planId,
		Environment: environment,
		Tags:        tags,
		NamespaceID: namespaceId,
		Endpoints:   []*Endpoint{},
		State:       StateStopped,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if _, err := services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO containers (id, name, image, plan_id, environment, tags, namespace_id, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, container.ID, container.Name, container.Image, container.PlanID, jsonEnvironment, jsonTags, container.NamespaceID,
		container.State, container.CreatedAt, container.UpdatedAt); err != nil {
		return nil, err
	}
	return container, nil
}

func UpdateContainer(ctx context.Context, container *Container) error {
	jsonEnvironment, err := json.Marshal(container.Environment)
	if err != nil {
		return err
	}
	jsonTags, err := json.Marshal(container.Tags)
	if err != nil {
		return err
	}
	_, err = services.Postgres().Client().ExecContext(ctx, `
		UPDATE containers SET name = $2, image = $3, plan_id = $4, environment = $5, tags = $6, namespace_id = $7,
		                      state = $8, created_at = $9, updated_at = $10
		WHERE id = $1
	`, container.ID, container.Name, container.Image, container.PlanID, string(jsonEnvironment), string(jsonTags),
		container.NamespaceID, container.State, container.CreatedAt, container.UpdatedAt)
	return err
}

func DeleteContainer(ctx context.Context, id string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		DELETE FROM containers WHERE id = $1
	`, id)
	return err
}

func FindContainerByID(ctx context.Context, id string) (*Container, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, image, plan_id, environment, tags, namespace_id, state, created_at, updated_at FROM containers
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return containerFromRows(rows)
	}
	return nil, nil
}

func FindContainerByNameAndNamespaceID(ctx context.Context, name, namespaceId string) (*Container, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, image, plan_id, environment, tags, namespace_id, state, created_at, updated_at FROM containers
		WHERE name = $1 AND namespace_id = $2
	`, name, namespaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return containerFromRows(rows)
	}
	return nil, nil
}

func FindContainersByNamespaceID(ctx context.Context, id string) ([]*Container, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, image, plan_id, environment, tags, namespace_id, state, created_at, updated_at FROM containers
		WHERE namespace_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var containers []*Container
	for rows.Next() {
		container, err := containerFromRows(rows)
		if err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}
	return containers, nil
}
