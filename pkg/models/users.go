package models

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	AvatarUrl string
	GithubID  int64
	Admin     bool
	CreatedAt time.Time
}

func userFromRows(rows *sql.Rows) (*User, error) {
	if rows.Next() {
		user := &User{}
		if err := rows.Scan(user.ID, user.Name, user.Email, user.AvatarUrl,
			user.GithubID, user.Admin, user.CreatedAt); err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

func UsersFindByID(ctx context.Context, id string) (*User, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, admin, created_at FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return userFromRows(rows)
}

func UsersFindByEmail(ctx context.Context, email string) (*User, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, admin, created_at FROM users
		WHERE email = $1
	`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return userFromRows(rows)
}

func UsersFindByGithubID(ctx context.Context, id int64) (*User, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, admin, created_at FROM users
		WHERE github_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return userFromRows(rows)
}

func UsersCreate(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, avatar_url, github_id, admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.Name, user.Email, user.AvatarUrl, user.GithubID, user.Admin, user.CreatedAt)
	return err
}

func UsersUpdate(ctx context.Context, user *User) error {
	_, err := db.ExecContext(ctx, `
		UPDATE users SET name = $2, email = $3, avatar_url = $4, github_id = $5, admin = $6, created_at = $7
	  	WHERE id = $1
	`, user.ID, user.Name, user.Email, user.AvatarUrl, user.GithubID, user.Admin, user.CreatedAt)
	return err
}

func UsersDelete(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM users WHERE id = $1
	`, id)
	return err
}
