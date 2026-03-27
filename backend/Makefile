.PHONY: build test lint migrate run

build:
	CGO_ENABLED=0 go build -trimpath -o bin/server ./cmd/server

run:
	GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server

migrate:
	GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only

test:
	GOTOOLCHAIN=local GOSUMDB=off go test ./...

lint:
	golangci-lint run ./...
