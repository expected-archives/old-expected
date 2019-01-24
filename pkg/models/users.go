package models

import (
	"context"
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

func UsersCreate(ctx context.Context, name, email, avatarUrl string, githubId int64, admin bool) (*User, error) {
	id := uuid.New().String()
	createdAt := time.Now()
	if _, err := db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, avatar_url, github_id, admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, name, email, avatarUrl, githubId, admin, createdAt); err != nil {
		return nil, err
	}
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		AvatarUrl: avatarUrl,
		GithubID:  githubId,
		Admin:     admin,
		CreatedAt: createdAt,
	}, nil
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

func UsersFindByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func UsersFindByGithubID(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}
