package models

import "context"

type User struct {
	ID        string
	Name      string
	Email     string
	AvatarUrl string
	GithubID  string
	Admin     bool
}

func UsersFindByID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func UsersFindByEmail(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func UsersFindByGithubID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}
