package containers

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/expectedsh/expected/pkg/services"
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
