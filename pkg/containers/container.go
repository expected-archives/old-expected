package containers

import (
	"database/sql"
	"time"
)

var db *sql.DB

type Container struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Endpoint    string            `json:"endpoint"`
	Memory      int               `json:"memory"`
	Environment map[string]string `json:"environment"`
	Tags        []string          `json:"tags"`
	OwnerID     string            `json:"-"`
	CreatedAt   time.Time         `json:"created_at"`
}

func InitDB(database *sql.DB) error {
	db = database
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS containers (
			id UUID NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			image VARCHAR(255) NOT NULL,
			endpoint VARCHAR(255) NOT NULL,
			memory INT NOT NULL,
			environment JSON NOT NULL,
			tags JSON NOT NULL,
			owner_id UUID NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	return err
}
