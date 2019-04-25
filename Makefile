.PHONY: test all check-counterfeiter gen-mocks

all: test

GO_SOURCES = $(shell find . -type f -name '*.go')

test:
	GO111MODULE=on go test ./... -coverprofile=coverage.txt -covermode=atomic

check-counterfeiter:
    # Use go get in GOPATH mode to install/update counterfeiter. This avoids polluting go.mod/go.sum.
	@which counterfeiter > /dev/null || (echo counterfeiter not found: issue "GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter" && false)

gen-mocks: check-counterfeiter
	counterfeiter -o pkg/registry/imagefakes/fake_image.go github.com/google/go-containerregistry/pkg/v1.Image
	counterfeiter -o pkg/registry/imagefakes/fake_image_index.go github.com/google/go-containerregistry/pkg/v1.ImageIndex
	counterfeiter pkg/registry LayoutPath

irel: $(GO_SOURCES)
	GO111MODULE=on go build -o irel cmd/irel/main.go

