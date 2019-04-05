package registrycli

import (
	"fmt"
	"net/http"
)

type DeleteStatus int

const (
	Deleted DeleteStatus = iota
	NotFound
	Unknown
)

func (d DeleteStatus) String() string {
	switch d {
	case Deleted:
		return "deleted"
	case NotFound:
		return "not found"
	}
	return "unknown"
}

func DeleteLayer(repo, digest string) (DeleteStatus, error) {
	token, err := newToken(repo)
	if err != nil {
		return Unknown, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/blobs/%s", registryUrl, repo, digest), nil)

	req.Header.Set("Authorization", "bearer "+token)
	res, err := client.Do(req)

	status := Unknown
	if res != nil && res.StatusCode >= 200 && res.StatusCode < 300 {
		status = Deleted
	} else if res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
		status = NotFound
	}

	return status, err
}
