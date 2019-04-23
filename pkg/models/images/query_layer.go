package images

import (
	"context"
	"database/sql"
	"errors"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/backoff"
	"time"
)

func layerFromRows(rows *sql.Rows) (*Layer, error) {
	layer := &Layer{}
	if err := rows.Scan(&layer.Digest, &layer.Repository, &layer.Size, &layer.CreatedAt, &layer.UpdatedAt); err != nil {
		return nil, err
	}
	return layer, nil
}

func CreateLayer(ctx context.Context, repo, digest string, size int64) (*Layer, error) {
	now := time.Now()
	_, err := services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO layers (repository, digest, size, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
	`, repo, digest, size)
	return &Layer{
		Repository: repo,
		Digest:     digest,
		Size:       size,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, err
}

func CreateLayers(ctx context.Context, layers []Layer) error {
	tx, err := services.Postgres().Client().Begin()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO layers (repository, digest, size, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (digest)
		DO UPDATE SET updated_at = now()
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, layer := range layers {
		if _, err := stmt.ExecContext(ctx, layer.Repository, layer.Digest, layer.Size); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func UpdateLayer(ctx context.Context, digest string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		UPDATE layers 
		SET updated_at = now()
		WHERE digest = $1
	`, digest)
	return err
}

func UpdateLayerRepository(ctx context.Context, digest string) error {
	_, e := services.Postgres().Client().ExecContext(ctx, `
		WITH res AS (
		    SELECT i.namespace_id as namespace, i.name as name
		    FROM public.image_layer
		             LEFT JOIN public.images i on image_layer.image_id = i.id
		    WHERE layer_digest = $1
		)
		UPDATE layers
		SET repository = concat((SELECT namespace FROM res), '/', (SELECT name FROM res)),
		    updated_at = now()
		WHERE digest = $1
		`, digest)
	return e
}

func IsLayerUnused(ctx context.Context, digest string) (bool, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT (SELECT count(*) FROM image_layer WHERE layer_digest = digest) FROM layers WHERE digest = $1
	`, digest)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		n := 0
		if err := rows.Scan(&n); err != nil {
			return false, err
		}
		return n == 0, nil
	}
	return false, errors.New("layer not found")
}

func DeleteLayerByDigest(ctx context.Context, digest string) error {
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
		DELETE FROM layers WHERE digest = $1
	`, digest)
	return err
}

func FindLayerByDigest(ctx context.Context, digest string) (*Layer, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT digest, repository, size, created_at, updated_at
		FROM layers
		WHERE digest = $1
	`, digest)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return layerFromRows(rows)
	}
	return nil, nil
}

func FindLayersByImageId(ctx context.Context, imageId string) ([]*Layer, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT
    		layers.digest,
    		layers.repository,
    		layers.size,
    		layers.created_at,
    		layers.updated_at
		FROM image_layer
    	LEFT JOIN layers ON image_layer.layer_digest = layers.digest
		WHERE image_id = $1
	`, imageId)
	if err != nil {
		return nil, err
	}
	var layers []*Layer
	for rows.Next() {
		if layer, err := layerFromRows(rows); err != nil {
			return nil, err
		} else {
			layers = append(layers, layer)
		}
	}
	return layers, nil
}

func FindUnusedLayers(ctx context.Context, olderThan time.Duration, limit int64) ([]*Layer, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT digest, repository, size, created_at, updated_at FROM layers 
		WHERE
			updated_at <  now() - interval '`+olderThan.String()+`' AND
			(SELECT count(*) FROM image_layer WHERE layer_digest = digest) = 0
		ORDER BY updated_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}

	var layers []*Layer
	for rows.Next() {
		if layer, err := layerFromRows(rows); err != nil {
			return nil, err
		} else {
			layers = append(layers, layer)
		}
	}
	return layers, nil
}
