package registrycli

import (
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/expectedsh/expected/pkg/registryserver/auth"
	"github.com/expectedsh/expected/pkg/registryserver/auth/token"
	"io/ioutil"
	"net/http"
	"time"
)

// todo change this in the future
const registryUrl = "http://localhost:5000"

// GetManifest return a manifest from a repo and digest.
func GetManifest(repo, digest string) *schema2.Manifest {
	manifest := &schema2.Manifest{}
	tok, _ := token.Generate(auth.RequestFromDaemon{
		Login:   "admin",
		Service: "registry",
	}, []auth.AuthorizedScope{
		{
			Scope: auth.Scope{
				Type: "repository",
				Name: repo,
			},
			AuthorizedActions: []string{"pull", "push", "delete"},
		},
	})
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v2/%s/manifests/%s", registryUrl, repo, digest), nil)

	req.Header.Set("Authorization", "bearer "+tok)
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

// getLayers call the registry to get all fs layers for a given digest and repo.
func GetLayers(repo, digest string, size int64) []images.Layer {
	manifest := GetManifest(repo, digest)

	if manifest == nil {
		return nil
	}

	var layers []images.Layer

	// add layer digest
	for _, layer := range manifest.Layers {
		layers = append(layers, images.Layer{
			Digest:    layer.Digest.String(),
			Size:      layer.Size,
			Count:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	layers = append(layers, images.Layer{
		Digest:    digest,
		Size:      size,
		Count:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	layers = append(layers, images.Layer{
		Digest:    manifest.Config.Digest.String(),
		Size:      manifest.Config.Size,
		Count:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return layers
}
