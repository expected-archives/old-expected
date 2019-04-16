package images

import (
	"context"
	"database/sql"
	"errors"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/backoff"
)

func statsFromRows(rows *sql.Rows) (*Stats, error) {
	stats := &Stats{}
	if err := rows.Scan(&stats.ImageID, &stats.NamespaceID, &stats.Digest, &stats.Name, &stats.Tag,
		&stats.Layers, &stats.Size); err != nil {
		return nil, err
	}
	return stats, nil
}

func CreateImageLayer(ctx context.Context, layers []Layer, imageId string) error {
	tx, err := services.Postgres().Client().Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO image_layer (image_id, layer_digest)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, layer := range layers {

		rows, err := tx.QueryContext(ctx, `
			SELECT * FROM image_layer WHERE image_id = $1 AND layer_digest = $2
		`, imageId, layer.Digest)
		if err != nil {
			return err
		}

		hasImageLayer := rows.Next()
		if err := rows.Close(); err != nil {
			return err
		}

		// insert into image_layer if there is no rows founded.
		if !hasImageLayer {
			if _, err := stmt.ExecContext(ctx, imageId, layer.Digest); err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		}

	}
	return tx.Commit()
}

func DeleteImageLayerByImageID(ctx context.Context, imageId string) error {
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
		DELETE FROM image_layer WHERE image_id = $1
	`, imageId)
	return err
}

func FindImagesStatsByNamespaceID(ctx context.Context, namespaceId string) ([]*Stats, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT
		    img.id,
			img.namespace_id,
       		img.digest,
       		img.name,
       		img.tag,
       		count(layers.created_at) AS layers,
       		sum(layers.size)         AS size
		FROM image_layer
       		LEFT JOIN layers ON image_layer.layer_digest = layers.digest
       		LEFT JOIN images img on image_layer.image_id = img.id
		WHERE namespace_id = $1
		GROUP BY img.name, img.digest, img.tag, img.namespace_id, img.id;
	`, namespaceId)
	if err != nil {
		return nil, err
	}
	var statsList []*Stats
	for rows.Next() {
		if stats, err := statsFromRows(rows); err != nil {
			return nil, err
		} else {
			statsList = append(statsList, stats)
		}
	}
	return statsList, nil
}

func FindLayerCountReferences(ctx context.Context, digest string) (int64, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT count(*) 
		FROM image_layer 
		WHERE layer_digest = $1
	`, digest)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		var count int64
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, errors.New("unexpected error with psql in FindLayerCountReferences function")
}
