package accounts

import (
	"context"
	"database/sql"
	"time"
)

var db *sql.DB

type Account struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	AvatarURL         string    `json:"avatar_url"`
	GithubID          int64     `json:"-"`
	GithubAccessToken string    `json:"-"`
	APIKey            string    `json:"api_key"`
	Admin             bool      `json:"admin"`
	CreatedAt         time.Time `json:"created_at"`
}

func InitDB(database *sql.DB) error {
	db = database
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id UUID NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			avatar_url VARCHAR(255) NOT NULL,
			github_id BIGINT NOT NULL,
			github_access_token VARCHAR(255) NOT NULL,
			api_key VARCHAR(32) NOT NULL,
			admin BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}
	return CreateAdmin(context.Background())
}
