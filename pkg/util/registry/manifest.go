package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// GetManifest return a manifest from a repo and digest.
func GetManifest(repo, digest string) *schema2.Manifest {
	manifest := &schema2.Manifest{}

	tok, _ := newToken(repo)

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
		layers = append(layers, images.Layer{Digest: layer.Digest.String(), Size: layer.Size})
	}

	layers = append(layers, images.Layer{Digest: digest, Size: size})
	layers = append(layers, images.Layer{Digest: manifest.Config.Digest.String(), Size: manifest.Config.Size})

	return layers
}

// Delete an image
// todo export it where the endpoint will be
func Delete(ctx context.Context, img images.Image) error {
	log := logrus.NewEntry(logrus.StandardLogger()).
		WithField("service", "registry-hook").
		WithField("event", "delete").
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag)

	log = log.WithField("id", img.ID).WithField("tag", img.Tag)
	log.Info()

	layers, err := images.FindLayersByImageId(ctx, img.ID)
	if err != nil {
		log.WithError(err).Error("finding layers by image id")
		return err
	}

	// deleting relations between image and layers
	if err := images.DeleteImageLayerByImageID(ctx, img.ID); err != nil {
		log.WithError(err).Error("deleting image_layer rows by image id")
		return err
	}

	for _, layer := range layers {

		// If layer is again referenced and unfortunately the repository property is the one that
		// the registry delete, another repository is choose.
		// Else the layer update_at property is updated to be garbage collected.

		layerLog := log.WithField("digest", layer.Digest)

		if cnt, err := images.FindLayerCountReferences(ctx, layer.Digest); err != nil {
			layerLog.WithError(err).Error("finding layer count references")
			return err
		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceID, img.Name) {
			if err := images.UpdateLayerRepository(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating repository of layer")
				return err
			}
		} else {
			if err := images.UpdateLayer(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating layer")
				return err
			}
		}
	}

	// deleting the image at the end to be sure all actions above has been executed
	if err := images.DeleteImage(ctx, img.ID); err != nil {
		log.WithError(err).Error("deleting image")
		return err
	}

	return nil
}
