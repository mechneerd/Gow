.PHONY: all test lint build coverage clean fmt vet mod

all: test lint build

test:
	go test ./... -race -v

lint:
	golangci-lint run ./...

build:
	@mkdir -p bin
	go build -o bin/gow ./cmd/gow
	go build -o bin/artisan ./cmd/artisan

coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean:
	rm -rf bin/ coverage.out coverage.html gow gow.exe artisan artisan.exe

fmt:
	go fmt ./...

vet:
	go vet ./...

mod:
	go mod tidy

# Install dev tools (run once)
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/cover@latest
