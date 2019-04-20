package images

import (
	"context"
	"github.com/expectedsh/expected/pkg/services"
)

func FindImagesName(ctx context.Context, namespaceId string) ([]string, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT concat(name, concat(':', tag)) 
		FROM images 
		WHERE namespace_id = $1 
		GROUP BY concat(name, concat(':', tag))
	`, namespaceId)
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		tag := ""
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		names = append(names, tag)
	}
	return names, nil
}
