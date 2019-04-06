package images

import (
	"context"
	"github.com/expectedsh/expected/pkg/util/backoff"
	"github.com/google/uuid"
	"time"
)

func Create(ctx context.Context, name, digest, namespaceId, tag string) (*Image, error) {
	id := uuid.New().String()
	now := time.Now()
	err := backoff.ExecContext(db, ctx, `
		INSERT INTO images (id, name, digest, namespace_id, tag, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, name, digest, namespaceId, tag, now)
	return &Image{
		ID:          id,
		NamespaceID: namespaceId,
		Digest:      digest,
		Name:        name,
		Tag:         tag,
		CreatedAt:   now,
	}, err
}

func CreateLayer(ctx context.Context, repo, digest string, size int64) (*Layer, error) {
	now := time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO layers (repository, digest, size)
		VALUES ($1, $2, $3)
	`, repo, digest, size)
	return &Layer{
		Repository: repo,
		Digest:     digest,
		Size:       size,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, err
}

func CreateLayers(ctx context.Context, layers []Layer, imageId string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO layers (repository, digest, size)
		VALUES ($1, $2, $3)
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

func CreateImageLayer(ctx context.Context, layers []Layer, imageId string) error {
	tx, err := db.Begin()
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

func UpdateLayer(ctx context.Context, digest string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE layers 
		SET updated_at = now()
		WHERE digest = $2
	`, digest)
	return err
}

func DeleteLayer(ctx context.Context, digest string) error {
	err := backoff.ExecContext(db, ctx, `
		DELETE FROM layers WHERE digest = $1
	`, digest)
	return err
}
