package models

import "context"

type User struct {
	ID        string
	Name      string
	Email     string
	AvatarUrl string
	GithubID  int64
	Admin     bool
}

func UsersCreate(ctx context.Context, id, name, email, avatarUrl string, githubId int64, admin bool) (*User, error) {
	return nil, nil
}

func UsersFindByID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func UsersFindByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func UsersFindByGithubID(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}
