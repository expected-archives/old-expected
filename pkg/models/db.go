package models

import (
	"database/sql"
	"time"
)

var (
	db       *sql.DB
	Accounts *AccountsModel
)

func InitDB(addr string, connMaxLifetime time.Duration, maxIdleConns, maxOpenConns int) error {
	var err error
	if db, err = sql.Open("postgres", addr); err != nil {
		return err
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
  			id UUID NOT NULL PRIMARY KEY,
  			name VARCHAR(255) NOT NULL,
  			email VARCHAR(255) NOT NULL,
  			avatar_url VARCHAR(255) NOT NULL,
  			github_id BIGINT NOT NULL,
  			github_access_token VARCHAR(255) NOT NULL,
  			github_refresh_token VARCHAR(255) NOT NULL,
  			admin BOOLEAN DEFAULT FALSE,
  			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	Accounts = &AccountsModel{}
	return err
}
