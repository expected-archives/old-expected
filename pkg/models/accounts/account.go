package accounts

import (
	"github.com/google/uuid"
	"strings"
	"time"
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

func (account *Account) RegenerateAPIKey() string {
	account.APIKey = strings.Replace(uuid.New().String(), "-", "", -1)
	return account.APIKey
}
