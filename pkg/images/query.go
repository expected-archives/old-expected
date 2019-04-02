package images

import "context"

func FindLayerByDigest(ctx context.Context, digest string) (*Layer, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT digest, count, size, created_at
		FROM layers
		WHERE digest = $1
	`, digest)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		layer := &Layer{}
		if err := rows.Scan(layer.Digest, layer.Count, layer.Size, layer.CreatedAt); err != nil {
			return nil, err
		}
		return layer, nil
	}
	return nil, nil
}

func LayerIncrement(ctx context.Context, digest string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE layers 
		SET count = count + 1 
		WHERE digest = $1
	`, digest)
	return err
}

func LayerDecrement(ctx context.Context, digest string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE layers 
		SET count = count - 1 
		WHERE digest = $1
	`, digest)
	return err
}

func FindImageByID(ctx context.Context, id string) (*Image, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, owner_id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		image := &Image{}
		if err := rows.Scan(image.ID, image.OwnerID, image.NamespaceID, image.Digest, image.Tag, image.Name,
			image.CreatedAt); err != nil {
			return nil, err
		}
		return image, nil
	}
	return nil, nil
}

func FindImagesByNamespaceID(ctx context.Context, namespaceId string) ([]*Image, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, owner_id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE namespace_id = $1
	`, namespaceId)
	if err != nil {
		return nil, err
	}
	var images []*Image
	for rows.Next() {
		image := &Image{}
		if err := rows.Scan(image.ID, image.OwnerID, image.NamespaceID, image.Digest, image.Tag, image.Name,
			image.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func FindImageByTagAndName(ctx context.Context, name, tag string) (*Image, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, owner_id, namespace_id, digest, tag, name, created_at
		FROM images
		WHERE name=$1 AND tag = $2
	`, name, tag)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		image := &Image{}
		if err := rows.Scan(image.ID, image.OwnerID, image.NamespaceID, image.Digest, image.Tag, image.Name,
			image.CreatedAt); err != nil {
			return nil, err
		}
		return image, nil
	}
	return nil, nil
}
