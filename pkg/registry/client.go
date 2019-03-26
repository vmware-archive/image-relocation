package registry

import (
	"github.com/pivotal/image-relocation/pkg/image"
)

// Client provides a way of interacting with image registries.
type Client interface {
	// Digest returns the digest of the given image or an error if the image does not exist or the digest is unavailable.
	Digest(image.Name) (image.Digest, error)

	// NewLayout creates a Layout for the Client and creates a corresponding directory containing a new OCI image layout at
	// the given file system path.
	NewLayout(path string) (Layout, error)

	// ReadLayout creates a Layout for the Client from the given file system path of a directory containing an existing
	// OCI image layout.
	ReadLayout(path string) (Layout, error)
}

type client struct {}

func NewRegistryClient() Client {
	return client{}
}

func (r client) Digest(n image.Name) (image.Digest, error) {
	img, err := readRemoteImage(n)
	if err != nil {
		return image.EmptyDigest, err
	}

	hash, err := img.Digest()
	if err != nil {
		return image.EmptyDigest, err
	}

	return image.NewDigest(hash.String()), nil
}