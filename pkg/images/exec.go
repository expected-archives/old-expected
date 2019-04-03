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

func CreateLayer(ctx context.Context, digest string, size int64) (*Layer, error) {
	now := time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO layers (digest, size, created_at)
		VALUES ($1, $2, $3)
	`, digest, size, now)
	return &Layer{
		Digest:    digest,
		Size:      size,
		CreatedAt: now,
		UpdatedAt: now,
	}, err
}

func CreateImageLayer(ctx context.Context, imageId, layerDigest string) (*ImageLayer, error) {
	now := time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO image_layer (image_id, layer_digest, created_at)
		VALUES ($1, $2, $3)
	`, imageId, layerDigest, now)
	return &ImageLayer{
		ImageID:     imageId,
		LayerDigest: layerDigest,
		CreatedAt:   now,
	}, err
}
