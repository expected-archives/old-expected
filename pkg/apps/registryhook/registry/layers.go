package registry

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"net/http"
	"time"
)

func DeleteLayer(repo, digest string) (DeleteStatus, error) {
	token, err := services.Auth().Client().GenerateToken(context.Background(), &protocol.GenerateTokenRequest{
		Image:    repo,
		Duration: int64(time.Minute * 10),
	})
	if err != nil {
		return DeleteStatusUnknown, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/blobs/%s", registryUrl, repo, digest), nil)

	req.Header.Set("Authorization", "bearer "+token.Token)
	res, err := client.Do(req)

	status := DeleteStatusUnknown
	if res != nil && res.StatusCode >= 200 && res.StatusCode < 300 {
		status = DeleteStatusDeleted
	} else if res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
		status = DeleteStatusNotFound
	}

	return status, err
}
