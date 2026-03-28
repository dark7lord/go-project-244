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

cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
# 	go tool cover -html=coverage.out


.PHONY: test fmt lint-fix lint build run