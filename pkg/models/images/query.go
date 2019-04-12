package images

import (
	"context"
	"database/sql"
	"errors"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/backoff"
	"github.com/google/uuid"
	"time"
)

func scanImage(rows *sql.Rows, image *Image) error {
	return rows.Scan(&image.ID, &image.NamespaceID, &image.Digest, &image.Tag, &image.Name, &image.CreatedAt, &image.DeleteMode)
}

func scanLayer(rows *sql.Rows, layer *Layer) error {
	return rows.Scan(&layer.Digest, &layer.Repository, &layer.Size, &layer.CreatedAt, &layer.UpdatedAt)
}

func scanStats(rows *sql.Rows, stats *Stats) error {
	return rows.Scan(&stats.ImageID, &stats.NamespaceID, &stats.Digest, &stats.Name, &stats.Tag, &stats.Layers, &stats.Size)
}

func Create(ctx context.Context, name, digest, namespaceId, tag string) (*Image, error) {
	id := uuid.New().String()
	now := time.Now()
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
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

func CreateImageLayerRelations(ctx context.Context, layers []Layer, imageId string) error {
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

func UpdateImageDeleteMode(ctx context.Context, imageId string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		UPDATE images 
		SET delete_mode = TRUE
		WHERE id = $1
	`, imageId)
	return err
}

func DeleteLayer(ctx context.Context, digest string) error {
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
		DELETE FROM layers WHERE digest = $1
	`, digest)
	return err
}

func DeleteImageLayerByImageID(ctx context.Context, imageId string) error {
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
		DELETE FROM image_layer WHERE image_id = $1
	`, imageId)
	return err
}

// DeleteImage delete an image only if there is no more relations
// in table image_layer.
func DeleteImage(ctx context.Context, imageId string) error {
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
		DELETE FROM images 
		WHERE id = $1 AND
		      (SELECT count(*) FROM image_layer WHERE image_id = id) = 0
	`, imageId)
	return err
}

// DeleteImageByDigest delete an image only if there is no more relations
// in table image_layer.
func DeleteImageByDigest(ctx context.Context, digest string) error {
	err := backoff.ExecContext(services.Postgres().Client(), ctx, `
		DELETE FROM images 
		WHERE digest = $1 AND 
		      (SELECT count(*) FROM image_layer WHERE image_id = id) = 0
	`, digest)
	return err
}

func FindLayerByDigest(ctx context.Context, digest string) (*Layer, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT *
		FROM layers
		WHERE digest = $1
	`, digest)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		layer := &Layer{}
		if err := scanLayer(rows, layer); err != nil {
			return nil, err
		}
		return layer, nil
	}
	return nil, nil
}

func FindImageByID(ctx context.Context, id string) (*Image, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at, delete_mode
		FROM images
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		image := &Image{}
		if err := scanImage(rows, image); err != nil {
			return nil, err
		}
		return image, nil
	}
	return nil, nil
}

func FindImageByInfos(ctx context.Context, namespaceId, name, tag, digest string) (*Image, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at, delete_mode
		FROM images
		WHERE namespace_id=$1 AND name=$2 AND tag=$3 AND digest = $4
	`, namespaceId, name, tag, digest)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		image := &Image{}
		if err := scanImage(rows, image); err != nil {
			return nil, err
		}
		return image, nil
	}
	return nil, nil
}

func FindLayersByImageId(ctx context.Context, imageId string) ([]*Layer, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT
    		layers.*
		FROM image_layer
    	LEFT JOIN layers ON image_layer.layer_digest = layers.digest
		WHERE image_id = $1
	`, imageId)
	if err != nil {
		return nil, err
	}
	var layers []*Layer
	for rows.Next() {
		layer := &Layer{}
		if err := scanLayer(rows, layer); err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	return layers, nil
}

func FindImageStatsByImageID(ctx context.Context, imageId string) (*Stats, error) {
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
		WHERE image_id = $1
		GROUP BY img.name, img.digest, img.tag, img.namespace_id, img.id;
	`, imageId)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		stats := &Stats{}
		if err := scanStats(rows, stats); err != nil {
			return nil, err
		}
		return stats, nil
	}
	return nil, nil
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
		stats := &Stats{}
		if err := scanStats(rows, stats); err != nil {
			return nil, err
		}
		statsList = append(statsList, stats)
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

func FindUnusedLayers(ctx context.Context, olderThan time.Duration, limit int64) ([]*Layer, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT * FROM layers 
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
		layer := &Layer{}
		if err := scanLayer(rows, layer); err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	return layers, nil
}
