package github

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io/ioutil"
)

type Config struct {
	ClientID     string `envconfig:"client_id"`
	ClientSecret string `envconfig:"client_secret"`
}

type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type Email struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

const githubURL = "https://api.github.com"

func GetUser(ctx context.Context, token *oauth2.Token) (*User, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(githubURL + "/user")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	if user.Login == "" || user.Name == "" || user.AvatarUrl == "" {
		return nil, errors.New("invalid github account")
	}
	return &user, nil
}

func GetEmails(ctx context.Context, token *oauth2.Token) (*[]Email, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(githubURL + "/user/emails")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var emails *[]Email
	if err = json.Unmarshal(body, &emails); err != nil {
		return nil, err
	}
	return emails, nil
}

func GetPrimaryEmail(ctx context.Context, token *oauth2.Token) (*Email, error) {
	emails, err := GetEmails(ctx, token)
	if err != nil {
		return nil, err
	}
	for _, email := range *emails {
		if email.Primary {
			return &email, nil
		}
	}
	return nil, nil
}
