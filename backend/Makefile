.PHONY: run build test lint

run:
	go run ./cmd/api

build:
	go build -o bin/server ./cmd/api

test:
	go test ./...

lint:
	docker run --rm -v "$(CURDIR)":/app -w /app golangci/golangci-lint:latest golangci-lint run
