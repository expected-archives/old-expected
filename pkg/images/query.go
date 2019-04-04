package images

import (
	"context"
	"database/sql"
	"github.com/expectedsh/expected/pkg/util/backoff"
	"strconv"
	"strings"
)

func scanImage(rows *sql.Rows, image *Image) error {
	return rows.Scan(&image.ID, &image.NamespaceID, &image.Digest, &image.Tag, &image.Name, &image.CreatedAt)
}

func scanLayer(rows *sql.Rows, layer *Layer) error {
	return rows.Scan(&layer.Digest, &layer.Size, &layer.Count, &layer.CreatedAt, &layer.UpdatedAt)
}

func scanStats(rows *sql.Rows, stats *Stats) error {
	return rows.Scan(&stats.NamespaceID, &stats.Digest, &stats.Name, &stats.Tag, &stats.Layers, &stats.Size)
}

// replaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func replaceSQL(old, searchPattern string) (string, int) {
	tmpCount := strings.Count(old, searchPattern)
	m := 1
	for ; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old, m
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

func InsertLayers(ctx context.Context, layers []Layer, imageId string) error {
	values := make([]interface{}, 0)

	query := `
		INSERT INTO layers (digest, size, count, created_at, updated_at)
		VALUES `

	for _, row := range layers {
		query += "(?, ?, ?, ?, ?), "
		values = append(values, row.Digest, row.Size, row.Count, row.CreatedAt, row.UpdatedAt)
	}

	query = query[0 : len(query)-2]
	query, n := replaceSQL(query, "?")
	query += `
		ON CONFLICT (digest)
		DO UPDATE
			SET updated_at = now(), 
				count = (SELECT count(*) 
					FROM image_layer 
					WHERE layer_digest=EXCLUDED.digest AND image_id<>$` + strconv.Itoa(n) + `)+1`

	values = append(values, imageId)
	return backoff.ExecContext(db, ctx, query, values...)
}

func InsertImageLayer(ctx context.Context, layers []Layer, imageId string) error {
	values := make([]interface{}, 0)

	query := `
		INSERT INTO image_layer (image_id, layer_digest)
		VALUES `

	for _, row := range layers {
		query += "(?, ?), "
		values = append(values, imageId, row.Digest)
	}

	query = query[0 : len(query)-2]
	query += `ON CONFLICT DO NOTHING`
	query, _ = replaceSQL(query, "?")
	return backoff.ExecContext(db, ctx, query, values...)
}

func LayerUpdate(ctx context.Context, digest, imageId string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE layers 
		SET updated_at = now(), 
			count = (	SELECT count(*) 
						FROM image_layer 
						WHERE layer_digest=EXCLUDED.digest AND image_id<>$1)+1
		WHERE digest = $2
	`, imageId, digest)
	return err
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

func FindImageByDigest(ctx context.Context, digest string) (*Image, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE digest = $1
	`, digest)
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

func FindImagesByNamespaceID(ctx context.Context, namespaceId string) ([]*Image, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE namespace_id = $1
	`, namespaceId)
	if err != nil {
		return nil, err
	}
	var images []*Image
	for rows.Next() {
		image := &Image{}
		if err := scanImage(rows, image); err == nil {
			images = append(images, image)
		}
	}
	return images, nil
}

func FindImageByInfos(ctx context.Context, name, tag, digest string) (*Image, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE name=$1 AND tag=$2 AND digest = $3
	`, name, tag, digest)
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

func GetStatsByImageId(ctx context.Context, imageId string) (*Stats, error) {
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
