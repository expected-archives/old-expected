package images

import (
	"database/sql"
	"time"
)

var db *sql.DB

type Image struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	OwnerID    string    `json:"-"`
	Digest     string    `json:"digest"`
	Repository string    `json:"repository"`
	Tag        string    `json:"tag"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}

func InitDB(database *sql.DB) error {
	db = database
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS images (
			id UUID NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			owner_id UUID NOT NULL,
			digest VARCHAR(255) NOT NULL,
			repository VARCHAR(255) NOT NULL,
			tag VARCHAR(255) NOT NULL,
			size INT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	return err
}
