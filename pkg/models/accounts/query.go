package accounts

import (
	"context"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/google/uuid"
	"time"
)

func Create(ctx context.Context, name, email, avatarUrl string, githubId int64,
	githubAccessToken string, admin bool) (*Account, error) {
	account := &Account{
		ID:                uuid.New().String(),
		Name:              name,
		Email:             email,
		AvatarURL:         avatarUrl,
		GithubID:          githubId,
		GithubAccessToken: githubAccessToken,
		Admin:             admin,
		CreatedAt:         time.Now(),
	}
	account.RegenerateAPIKey()
	_, err := services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO accounts (id, name, email, avatar_url, github_id, github_access_token, api_key, admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, account.ID, account.Name, account.Email, account.AvatarURL, account.GithubID, account.GithubAccessToken,
		account.APIKey, account.Admin, account.CreatedAt)
	return account, err
}

func Update(ctx context.Context, account *Account) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		UPDATE accounts SET name = $2, email = $3, avatar_url = $4, github_id = $5, github_access_token = $6,
		api_key = $7, admin = $8 WHERE id = $1
	`, account.ID, account.Name, account.Email, account.AvatarURL, account.GithubID, account.GithubAccessToken,
		account.APIKey, account.Admin)

	return err
}

func Delete(ctx context.Context, id string) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		DELETE FROM accounts WHERE id = $1
	`, id)

	return err
}

func FindByID(ctx context.Context, id string) (*Account, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, github_access_token, api_key, admin, created_at
		FROM accounts WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.AvatarURL, &account.GithubID,
			&account.GithubAccessToken, &account.APIKey, &account.Admin, &account.CreatedAt); err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}

func FindByGithubID(ctx context.Context, id int64) (*Account, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, github_access_token, api_key, admin, created_at
		FROM accounts WHERE github_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.AvatarURL, &account.GithubID,
			&account.GithubAccessToken, &account.APIKey, &account.Admin, &account.CreatedAt); err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}

func FindByAPIKey(ctx context.Context, apiKey string) (*Account, error) {
	rows, err := services.Postgres().Client().QueryContext(ctx, `
		SELECT id, name, email, avatar_url, github_id, github_access_token, api_key, admin, created_at
		FROM accounts WHERE api_key = $1
	`, apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.AvatarURL, &account.GithubID,
			&account.GithubAccessToken, &account.APIKey, &account.Admin, &account.CreatedAt); err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}
