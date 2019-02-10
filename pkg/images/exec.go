package images

import (
	"context"
	"github.com/google/uuid"
	"time"
)

func Create(ctx context.Context, name, ownerId, digest, repository, tag string, size int64) (*Image, error) {
	id := uuid.New().String()
	now := time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO images (id, name, owner_id, digest, repository, tag, size, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, id, name, ownerId, digest, repository, tag, size, now)
	return &Image{
		ID:         id,
		Name:       name,
		OwnerID:    ownerId,
		Digest:     digest,
		Repository: repository,
		Tag:        tag,
		Size:       size,
		CreatedAt:  now,
	}, err
}
