package accounts

import "context"

func FindByID(ctx context.Context, id string) (*Account, error) {
	rows, err := db.QueryContext(ctx, `
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
	rows, err := db.QueryContext(ctx, `
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
	rows, err := db.QueryContext(ctx, `
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
