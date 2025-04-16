.PHONY: build run install clean

# Build the application
build:
	go build -o gocommit

# Run the application
run: build
	./gocommit

# Install the application
install:
	go install

# Clean build artifacts
clean:
	rm -f gocommit

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download 