.PHONY: build test clean install proto

# Build brisk CLI
build:
	go build -o bin/brisk ./cmd/brisk

# Build brisk server
build-server:
	go build -o bin/brisk-server ./cmd/brisk-server

# Build all
build-all: build build-server

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install brisk CLI
install:
	go install ./cmd/brisk

# Install brisk server
install-server:
	go install ./cmd/brisk-server

# Install all
install-all: install install-server

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/api/cache.proto

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Run brisk example
example:
	@echo "Running example build with brisk..."
	./bin/brisk run go-build
