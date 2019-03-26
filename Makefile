.PHONY: test all check-counterfeiter gen-mocks

all: test

test:
	GO111MODULE=on go test ./...

check-counterfeiter:
    # Use go get in GOPATH mode to install/update counterfeiter. This avoids polluting go.mod/go.sum.
	@which counterfeiter > /dev/null || (echo counterfeiter not found: issue "GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter" && false)

gen-mocks: check-counterfeiter
	counterfeiter pkg/registry Layout
	counterfeiter pkg/registry Client
