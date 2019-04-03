package images

import (
	"database/sql"
	"time"
)

var db *sql.DB

// Image is represented by a Manifest in the registry.
type Image struct {
	ID          string    `json:"id"`           // randomly generated uuid
	NamespaceID string    `json:"namespace_id"` // uuid of the namespace
	Digest      string    `json:"digest"`       // digest of the image (a sha256)
	Tag         string    `json:"tag"`          // tag version name: latest, v1, v2...
	Name        string    `json:"name"`         // name of the tag
	CreatedAt   time.Time `json:"created_at"`   // when the image was created
}

// Layer is used by Image.
// An image can contain X layers and a layer can be used by Y images.
type Layer struct {
	Digest    string    `json:"digest"`     // digest sha256 id of the layer
	Size      int64     `json:"size"`       // size of the layer in bytes
	Count     uint64    `json:"-"`          // count the number of image that use this layer
	CreatedAt time.Time `json:"created_at"` // when the layer was first registered
	UpdatedAt time.Time `json:"updated_at"` // the last time the layer was updated
}

// ImageLayer represent the relation between an Image and X Layer.
// This table will be used to count total size of an image, number of layers ...
type ImageLayer struct {
	ImageID     string    `json:"image_id"`     // image id that refer to the Image
	LayerDigest string    `json:"layer_digest"` // layer digest that refer to the Layer
	CreatedAt   time.Time `json:"created_at"`   // date the image id was linked to the layer digest
}

// Stats get statistic about an image.
type Stats struct {
	NamespaceID string    `json:"namespace_id"` // uuid of the namespace
	Digest      string    `json:"digest"`       // digest of the image
	Name        string    `json:"name"`         // name of the image
	Tag         string    `json:"tag"`          // name of the tag
	Layers      int       `json:"layers"`       // number of layers
	Size        int64     `json:"size"`         // total size of the image in bytes
	CreatedAt   time.Time `json:"created_at"`   // creation date of the image
}

func InitDB(database *sql.DB) error {
	db = database
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id UUID NOT NULL PRIMARY KEY,
			namespace_id UUID NOT NULL,
			digest TEXT NOT NULL,
			name VARCHAR(255) NOT NULL,
			tag VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);
	`)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS layers (
			digest TEXT NOT NULL PRIMARY KEY,
			size BIGINT NOT NULL,
			count BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now()
		);
	`)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS image_layer (
			image_id UUID NOT NULL REFERENCES images (id),
			layer_digest TEXT NOT NULL REFERENCES layers (digest),
			created_at TIMESTAMP DEFAULT now()
		);
	`)
	return err
}
