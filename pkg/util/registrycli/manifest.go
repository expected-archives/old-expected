package registrycli

import (
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/expectedsh/expected/pkg/registryserver/auth"
	"github.com/expectedsh/expected/pkg/registryserver/auth/token"
	"io/ioutil"
	"net/http"
)

func GetManifest(registryUrl, email, repo, digest string) *schema2.Manifest {
	manifest := &schema2.Manifest{}
	tok, _ := token.Generate(auth.RequestFromDaemon{
		Login:   email,
		Service: "registry",
	}, []auth.AuthorizedScope{
		{
			Scope: auth.Scope{
				Type: "repository",
				Name: repo,
			},
			AuthorizedActions: []string{"pull"},
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
	fmt.Println("bytes", string(bytes))
	if err != nil {
		return nil
	}
	fmt.Println("manifest", manifest)
	return manifest
}
