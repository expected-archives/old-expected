package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// GetManifest return a manifest from a repo and digest.
func GetManifest(repo, digest string) *schema2.Manifest {
	manifest := &schema2.Manifest{}

	token, err := services.Auth().Client().GenerateToken(context.Background(), &protocol.GenerateTokenRequest{
		Image:    repo,
		Duration: int64(time.Minute * 10),
		Scopes:   []protocol.Scope{protocol.Scope_PULL},
	})

	if err != nil {
		logrus.WithError(err).Error("can't generate token")
		return nil
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v2/%s/manifests/%s", registryUrl, repo, digest), nil)

	req.Header.Set("Authorization", "bearer "+token.Token)
	res, err := client.Do(req)

	if err != nil {
		return nil
	}

	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, manifest)
	if err != nil {
		return nil
	}
	return manifest
}

func DeleteManifest(repo, digest string) (DeleteStatus, error) {
	token, err := services.Auth().Client().GenerateToken(context.Background(), &protocol.GenerateTokenRequest{
		Image:    repo,
		Duration: int64(time.Minute * 10),
		Scopes:   []protocol.Scope{protocol.Scope_PULL, protocol.Scope_DELETE},
	})
	if err != nil {
		return DeleteStatusUnknown, err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/manifests/%s", registryUrl, repo, digest), nil)

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

// getLayers call the registry to get all fs layers for a given digest and repo.
func GetLayers(repo, digest string, size int64) []images.Layer {
	manifest := GetManifest(repo, digest)

	if manifest == nil {
		return nil
	}

	var layers []images.Layer

	// add layer digest
	for _, layer := range manifest.Layers {
		layers = append(layers, images.Layer{Digest: layer.Digest.String(), Size: layer.Size})
	}

	layers = append(layers, images.Layer{Digest: digest, Size: size})
	layers = append(layers, images.Layer{Digest: manifest.Config.Digest.String(), Size: manifest.Config.Size})

	return layers
}
