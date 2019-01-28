package models

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type AccountsModel struct{}

type Account struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	AvatarUrl          string    `json:"avatar_url"`
	GithubID           int64     `json:"-"`
	GithubAccessToken  string    `json:"-"`
	GithubRefreshToken string    `json:"-"`
	Admin              bool      `json:"admin"`
	CreatedAt          time.Time `json:"created_at"`
}

func (m AccountsModel) GetByID(ctx context.Context, id string) (*Account, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, github_access_token, github_refresh_token,  admin, created_at
		FROM accounts WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account := &Account{}
		if err := rows.Scan(account.ID, account.Name, account.Email, account.AvatarUrl, account.GithubID,
			account.GithubAccessToken, account.GithubRefreshToken, account.Admin, account.CreatedAt); err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}

func (m AccountsModel) GetByGithubID(ctx context.Context, id int64) (*Account, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, github_access_token, github_refresh_token, admin, created_at
		FROM accounts WHERE github_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account := &Account{}
		if err := rows.Scan(account.ID, account.Name, account.Email, account.AvatarUrl, account.GithubID,
			account.GithubAccessToken, account.GithubRefreshToken, account.Admin, account.CreatedAt); err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}

func (m AccountsModel) Create(ctx context.Context, account *Account) error {
	account.ID = uuid.New().String()
	account.CreatedAt = time.Now()
	_, err := db.ExecContext(ctx, `
		INSERT INTO accounts (id, name, email, avatar_url, github_id, github_access_token, github_refresh_token, admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, account.ID, account.Name, account.Email, account.AvatarUrl, account.GithubID,
		account.GithubAccessToken, account.GithubRefreshToken, account.Admin, account.CreatedAt)
	return err
}

func (m AccountsModel) Update(ctx context.Context, account *Account) error {
	_, err := db.ExecContext(ctx, `
		UPDATE accounts SET name = $2, email = $3, avatar_url = $4, github_id = $5, github_access_token = $6,
		github_refresh_token = $7, admin = $7 WHERE id = $1
	`, account.ID, account.Name, account.Email, account.AvatarUrl, account.GithubID, account.GithubAccessToken,
		account.GithubRefreshToken, account.Admin)
	return err
}

func (m AccountsModel) Delete(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM accounts WHERE id = $1
	`, id)
	return err
}
