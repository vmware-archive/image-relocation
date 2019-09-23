module github.com/pivotal/image-relocation

require (
	github.com/docker/distribution v2.7.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/google/go-containerregistry v0.0.0-20190729175742-ef12d49c8daf
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.3
	gomodules.xyz/jsonpatch/v2 v2.0.1
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.2
)

go 1.12
