package images

import (
	"context"
	"github.com/google/uuid"
	"time"
)

func Create(ctx context.Context, name, ownerId, digest, namespaceId, tag string) (*Image, error) {
	id := uuid.New().String()
	now := time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO images (id, name, owner_id, digest, namespace_id, tag, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, name, ownerId, digest, namespaceId, tag, now)
	return &Image{
		ID:          id,
		OwnerID:     ownerId,
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
