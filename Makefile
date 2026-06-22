BIN := bin/gendiff

build:
	go build -o $(BIN) ./cmd/gendiff/main.go

install:
	go install ./cmd/gendiff

uninstall:
	go clean -i ./cmd/gendiff

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

fmt:
	golangci-lint fmt

test:
	go test ./... -v


ARGS ?= file1.json file2.json

run: build
	$(BIN) $(ARGS)

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

cover-html: cover
	go tool cover -html=coverage.out

check: test lint build clean
	@echo "✅ All checks passed"

clean:
	rm -f $(BIN)
	rm -f coverage.out

.PHONY: install uninstall test fmt lint-fix lint build cover cover-html clean run check 