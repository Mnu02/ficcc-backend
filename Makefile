.PHONY: run test build clean fmt lint

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Build the application
build:
	go build -o bin/ficcc-backend main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Run linter
lint:
	golangci-lint run

# Display test coverage
coverage: test
	go tool cover -html=coverage.out
