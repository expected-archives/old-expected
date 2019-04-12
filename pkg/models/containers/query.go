package containers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"strings"
	"time"
)

func containerFromRows(rows *sql.Rows) (*Container, error) {
	var environmentJson, tagsJson string
	container := &Container{}
	err := rows.Scan(&container.ID, &container.Name, &container.Image, &container.Endpoint, &container.Memory,
		&environmentJson, &tagsJson, &container.OwnerID, &container.CreatedAt)
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

	return container, nil
}

func Create(ctx context.Context, name, image string, memory int, environment map[string]string,
	tags []string, ownerId string) (*Container, error) {
	id := uuid.New().String()
	endpoint := strings.Replace(id, "-", "", -1) + ".ctr.expected.sh"
	createdAt := time.Now()
	jsonEnvironment, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	_, err = services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO containers (id, name, image, endpoint, memory, environment, tags, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, id, name, image, endpoint, memory, string(jsonEnvironment), string(jsonTags), ownerId, createdAt)

	return &Container{
		ID:          id,
		Name:        name,
		Image:       image,
		Endpoint:    endpoint,
		Memory:      memory,
		Environment: environment,
		Tags:        tags,
		OwnerID:     ownerId,
		CreatedAt:   createdAt,
	}, err
}

func Update(ctx context.Context, container *Container) error {
	jsonEnvironment, err := json.Marshal(container.Environment)
	if err != nil {
		return err
	}
	jsonTags, err := json.Marshal(container.Tags)
	if err != nil {
		return err
	}

	_, err = services.Postgres().Client().ExecContext(ctx, `
		UPDATE containers SET name = $2, image = $3, endpoint = $4, memory = $5, environment = $6, tags = $7
		WHERE id = $1
	`, container.ID, container.Name, container.Image, container.Endpoint, container.Memory, string(jsonEnvironment),
		string(jsonTags))

	return err
}

func Delete(ctx context.Context, id string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		DELETE FROM containers WHERE id = $1
	`, id)

	return err
}

func FindByID(ctx context.Context, id string) (*Container, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, image, endpoint, memory, environment, tags, owner_id, created_at FROM containers
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

func FindByOwnerID(ctx context.Context, id string) ([]*Container, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, image, endpoint, memory, environment, tags, owner_id, created_at FROM containers
		WHERE owner_id = $1
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
