package accounts

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

func Create(ctx context.Context, name, email, avatarUrl string, githubId int64,
	githubAccessToken string, admin bool) (*Account, error) {
	id := uuid.New().String()
	apiKey := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now()

	_, err := db.ExecContext(ctx, `
		INSERT INTO accounts (id, name, email, avatar_url, github_id, github_access_token, api_key, admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, id, name, email, avatarUrl, githubId, githubAccessToken, apiKey, admin, createdAt)

	return &Account{
		ID:                id,
		Name:              name,
		Email:             email,
		AvatarURL:         avatarUrl,
		GithubID:          githubId,
		GithubAccessToken: githubAccessToken,
		APIKey:            apiKey,
		Admin:             admin,
		CreatedAt:         createdAt,
	}, err
}

func Update(ctx context.Context, account *Account) error {
	_, err := db.ExecContext(ctx, `
		UPDATE accounts SET name = $2, email = $3, avatar_url = $4, github_id = $5, github_access_token = $6,
		api_key = $7, admin = $8 WHERE id = $1
	`, account.ID, account.Name, account.Email, account.AvatarURL, account.GithubID, account.GithubAccessToken,
		account.APIKey, account.Admin)

	return err
}

func Delete(ctx context.Context, id string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM accounts WHERE id = $1
	`, id)

	return err
}
