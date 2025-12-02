VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go mod tidy
	go fmt
	go test
	go build \
		-ldflags "-X main.version=$(VERSION)" \
		-o ./bin/bedlamb
