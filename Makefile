BIN := bin/gendiff

build:
	go build -o $(BIN) ./cmd/gendiff/main.go

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
	go tool cover -html=coverage.out

clean:
	rm -f $(BIN)


.PHONY: test fmt lint-fix lint build cover clean run