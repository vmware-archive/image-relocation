.PHONY: test all

all: test

test:
	GO111MODULE=on go test ./... -coverprofile=coverage.txt -covermode=atomic
