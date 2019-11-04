.PHONY: test all check-counterfeiter gen-mocks release

all: test

OUTPUT = ./irel
GO_SOURCES = $(shell find . -type f -name '*.go')
VERSION ?= $(shell cat VERSION)
GITSHA = $(shell git rev-parse HEAD)
GITDIRTY = $(shell git diff --quiet HEAD || echo "dirty")
LDFLAGS_VERSION = -X github.com/pivotal/image-relocation/pkg/irel.cli_version=$(VERSION) \
				  -X github.com/pivotal/image-relocation/pkg/irel.cli_gitsha=$(GITSHA) \
				  -X github.com/pivotal/image-relocation/pkg/irel.cli_gitdirty=$(GITDIRTY)

test:
	GO111MODULE=on go test ./... -coverprofile=coverage.txt -covermode=atomic

check-counterfeiter:
    # Use go get in GOPATH mode to install/update counterfeiter. This avoids polluting go.mod/go.sum.
	@which counterfeiter > /dev/null || (echo counterfeiter not found: issue "GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter" && false)

gen-mocks: check-counterfeiter
	counterfeiter -o pkg/registry/ggcrfakes/fake_image.go github.com/google/go-containerregistry/pkg/v1.Image
	counterfeiter -o pkg/registry/ggcrfakes/fake_image_index.go github.com/google/go-containerregistry/pkg/v1.ImageIndex
	counterfeiter pkg/registry/ggcr/path LayoutPath
	counterfeiter pkg/registry Image
	counterfeiter -o pkg/registry/ggcr/registryclientfakes/fake_registry_client.go ./pkg/registry/ggcr RegistryClient

irel: $(GO_SOURCES)
	GO111MODULE=on go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT) cmd/irel/main.go

release: $(GO_SOURCES) test
	GOOS=darwin   GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     cmd/irel/main.go && tar -czf irel-darwin-amd64.tgz  $(OUTPUT)     && rm -f $(OUTPUT)
	GOOS=linux    GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     cmd/irel/main.go && tar -czf irel-linux-amd64.tgz   $(OUTPUT)     && rm -f $(OUTPUT)
	GOOS=windows  GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT).exe cmd/irel/main.go && zip -mq  irel-windows-amd64.zip $(OUTPUT).exe && rm -f $(OUTPUT).exe
