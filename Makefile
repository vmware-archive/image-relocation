.PHONY: test all

all: test

test:
	GO111MODULE=on go test ./...
