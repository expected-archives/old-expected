package images

import (
	"context"
	"database/sql"
	"time"
)

func scanImage(rows *sql.Rows, image *Image) error {
	return rows.Scan(&image.ID, &image.NamespaceID, &image.Digest, &image.Tag, &image.Name, &image.CreatedAt)
}

func scanLayer(rows *sql.Rows, layer *Layer) error {
	return rows.Scan(&layer.Digest, &layer.Repository, &layer.Size, &layer.CreatedAt, &layer.UpdatedAt)
}

func scanStats(rows *sql.Rows, stats *Stats) error {
	return rows.Scan(&stats.NamespaceID, &stats.Digest, &stats.Name, &stats.Tag, &stats.Layers, &stats.Size)
}

func FindLayerByDigest(ctx context.Context, digest string) (*Layer, error) {
	rows, err := db.QueryContext(ctx, `
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
	rows, err := db.QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
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
	rows, err := db.QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
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

func FindImageStatsByImageID(ctx context.Context, imageId string) (*Stats, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT
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
		GROUP BY img.name, img.digest, img.tag, img.namespace_id;
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
	rows, err := db.QueryContext(ctx, `
		SELECT
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
		GROUP BY img.name, img.digest, img.tag, img.namespace_id;
	`, namespaceId)
	if err != nil {
		return nil, err
	}
	var statsList []*Stats
	if rows.Next() {
		stats := &Stats{}
		if err := scanStats(rows, stats); err != nil {
			return nil, err
		}
		statsList = append(statsList, stats)
	}
	return statsList, nil
}

func FindUnusedLayers(ctx context.Context, olderThan time.Duration, limit int64) ([]*Layer, error) {
	rows, err := db.QueryContext(ctx, `
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
