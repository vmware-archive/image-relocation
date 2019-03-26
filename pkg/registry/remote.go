package registry

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pivotal/image-relocation/pkg/image"
	"net/http"
)

func readRemoteImage(n image.Name) (v1.Image, error) {
	auth, err := resolve(n)
	if err != nil {
		return nil, err
	}

	ref, err := name.ParseReference(n.String(), name.StrictValidation)
	if err != nil {
		return nil, err
	}

	return remote.Image(ref, remote.WithAuth(auth))
}

func writeRemoteImage(i v1.Image, n image.Name) error {
	auth, err := resolve(n)
	if err != nil {
		return err
	}

	ref, err := name.ParseReference(n.String(), name.WeakValidation)
	if err != nil {
		return err
	}

	return remote.Write(ref, i, auth, http.DefaultTransport)
}

func resolve(n image.Name) (authn.Authenticator, error) {
	repo, err := name.NewRepository(n.WithoutTag().String(), name.WeakValidation)
	if err != nil {
		return nil, err
	}

	return authn.DefaultKeychain.Resolve(repo.Registry)
}