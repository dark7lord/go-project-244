build:
	go build -o bin/gendiff ./cmd/gendiff/main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

fmt:
	golangci-lint fmt

test:
	go test

run:
	@go run ./cmd/gendiff/main.go

.PHONY: test fmt lint-fix lint build run