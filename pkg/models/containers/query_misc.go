package containers

import (
	"context"
	"github.com/expectedsh/expected/pkg/services"
)

func FindTagsByNamespaceID(ctx context.Context, id string) ([]string, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT ('[' || json_array_elements(tags)::text || ']')::json ->> 0 AS tag 
		FROM containers 
		WHERE namespace_id = $1 GROUP BY tag
	`, id)
	if err != nil {
		return nil, err
	}
	var tags []string
	for rows.Next() {
		tag := ""
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
