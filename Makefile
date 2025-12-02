BINARY=creality2mqtt

.PHONY: all build test fmt vet check release

all: build

build:
	go build -o $(BINARY) ./cmd/creality2mqtt

fmt:
	gofmt -s -w .

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet test

release:
	goreleaser release --clean
