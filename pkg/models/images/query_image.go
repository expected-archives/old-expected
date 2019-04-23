package images

import (
	"context"
	"database/sql"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/backoff"
	"github.com/google/uuid"
	"time"
)

func imageFromRows(rows *sql.Rows) (*Image, error) {
	image := &Image{}
	if err := rows.Scan(&image.ID, &image.NamespaceID, &image.Digest, &image.Tag, &image.Name, &image.CreatedAt); err != nil {
		return nil, err
	}
	return image, nil
}

func CreateImage(ctx context.Context, name, digest, namespaceId, tag string) (*Image, error) {
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

// DeleteImageByID delete an image only if there is no more relations
// in table image_layer.
func DeleteImageByID(ctx context.Context, imageId string) error {
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

func FindImageByID(ctx context.Context, id string) (*Image, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return imageFromRows(rows)
	}
	return nil, nil
}

func FindImageByInfos(ctx context.Context, namespaceId, name, tag, digest string) (*Image, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE namespace_id=$1 AND name=$2 AND tag=$3 AND digest = $4
	`, namespaceId, name, tag, digest)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return imageFromRows(rows)
	}
	return nil, nil
}

func FindImageDetail(ctx context.Context, namespaceId, name, tag string) (*ImageDetail, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE namespace_id=$1 AND name=$2 AND tag=$3
		ORDER BY created_at DESC
	`, namespaceId, name, tag)
	if err != nil {
		return nil, err
	}
	var manifests []Manifest
	for rows.Next() {
		if img, err := imageFromRows(rows); err != nil {
			return nil, err
		} else {
			if layers, err := FindLayersByImageId(ctx, img.ID); err != nil {
				return nil, err
			} else {
				manifests = append(manifests, Manifest{Image: img, Layers: layers})
			}
		}
	}
	if len(manifests) == 0 {
		return nil, nil
	} else {
		return &ImageDetail{
			ImageSummary: &ImageSummary{
				NamespaceID: manifests[0].Image.NamespaceID,
				Name:        manifests[0].Image.Name,
				Tag:         manifests[0].Image.Tag,
				LastPushAt:  manifests[0].Image.CreatedAt,
			},
			Manifests: manifests,
		}, nil
	}
}
