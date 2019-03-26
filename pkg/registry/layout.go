package registry

import (
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/pivotal/image-relocation/pkg/image"
	"os"
)

const (
	outputDirPermissions = 0755
)

// A Layout allows a registry client to interact with an OCI image layout on disk.
type Layout interface {
	// Add adds the image at the given image reference to the layout and returns the image's digest.
	Add(n image.Name) (image.Digest, error)

	// Push pushes the image with the given digest from the layout to the given image reference.
	Push(digest image.Digest, name image.Name) error
}

func (r client) NewLayout(path string) (Layout, error) {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.MkdirAll(path, outputDirPermissions); err != nil {
			return nil, err
		}
	}

	lp, err := layout.Write(path, empty.Index)
	if err != nil {
		return nil, err
	}

	return &imageLayout{
		registryClient: r,
		layoutPath: lp,
	}, nil
}

func (r client) ReadLayout(path string) (Layout, error) {
	lp, err := layout.FromPath(path)
	if err != nil {
		return nil, err
	}
	return &imageLayout{
		registryClient: r,
		layoutPath:     lp,
	}, nil
}

type imageLayout struct {
	registryClient client
	layoutPath     layout.Path
}

func (l *imageLayout) Add(n image.Name) (image.Digest, error) {
	img, err := readRemoteImage(n)
	if err != nil {
		return image.EmptyDigest, err
	}

	annotations := map[string]string{
		"org.opencontainers.image.ref.name": n.String(),
	}
	if err := l.layoutPath.AppendImage(img, layout.WithAnnotations(annotations)); err != nil {
		return image.EmptyDigest, err
	}

	hash, err := img.Digest()
	if err != nil {
		return image.EmptyDigest, err
	}

	return image.NewDigest(hash.String()), nil
}

func (l *imageLayout) Push(digest image.Digest, n image.Name) error {
	hash, err := v1.NewHash(digest.String())
	if err != nil {
		return err
	}
	imageIndex, err := l.layoutPath.ImageIndex()
	if err != nil {
		return err
	}
	i, err := imageIndex.Image(hash)
	if err != nil {
		return err
	}

	return writeRemoteImage(i, n)
}
