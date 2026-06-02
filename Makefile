.PHONY: run stop test build clean fmt lint

# Run the application (frees port 8080 first so a stale server can't linger
# with an out-of-date .env)
run:
	@PID=$$(lsof -tiTCP:8080 -sTCP:LISTEN 2>/dev/null); \
	if [ -n "$$PID" ]; then echo "Stopping existing server on :8080 (PID $$PID)"; kill $$PID; sleep 1; fi
	go run .

# Stop any server currently holding port 8080
stop:
	@PID=$$(lsof -tiTCP:8080 -sTCP:LISTEN 2>/dev/null); \
	if [ -n "$$PID" ]; then echo "Stopping server on :8080 (PID $$PID)"; kill $$PID; else echo "Nothing listening on :8080"; fi

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Build the application
build:
	go build -o bin/ficcc-backend .

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
