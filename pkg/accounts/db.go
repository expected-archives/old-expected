package accounts

import "database/sql"

var db *sql.DB

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
	return err
}
