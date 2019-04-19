package images

import (
	"context"
	"database/sql"
	"errors"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/backoff"
)

func imageSummaryFromRows(rows *sql.Rows) (*ImageSummary, error) {
	summary := &ImageSummary{}
	if err := rows.Scan(&summary.NamespaceID, &summary.Name, &summary.Tag, &summary.LastPushAt); err != nil {
		return nil, err
	}
	return summary, nil
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
			SELECT image_id, layer_digest, created_at
			FROM image_layer WHERE image_id = $1 AND layer_digest = $2
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

func FindImagesSummariesByNamespaceID(ctx context.Context, namespaceId string) ([]*ImageSummary, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT 
		    img.namespace_id,
      		img.name,
      		img.tag,
      		(
      		    SELECT layers.created_at
      		    FROM layers
      		    WHERE repository = concat(img.namespace_id, concat('/', img.name))
      		    ORDER BY layers.created_at DESC
      		    LIMIT 1
      		) AS last_push
		FROM image_layer
		         LEFT OUTER JOIN images img on image_layer.image_id = img.id
		WHERE namespace_id = $1
		GROUP BY img.name, img.tag, img.namespace_id;
	`, namespaceId)
	if err != nil {
		return nil, err
	}
	var statsList []*ImageSummary
	for rows.Next() {
		if stats, err := imageSummaryFromRows(rows); err != nil {
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
