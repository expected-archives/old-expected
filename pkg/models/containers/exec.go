package containers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"
)

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

	_, err = db.ExecContext(ctx, `
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

	_, err = db.ExecContext(ctx, `
		UPDATE containers SET name = $2, image = $3, endpoint = $4, memory = $5, environment = $6, tags = $7
		WHERE id = $1
	`, container.ID, container.Name, container.Image, container.Endpoint, container.Memory, string(jsonEnvironment),
		string(jsonTags))

	return err
}

func Delete(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM containers WHERE id = $1
	`, id)

	return err
}
