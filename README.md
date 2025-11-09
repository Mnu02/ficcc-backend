# ficcc-backend

A Golang backend service for FICCC.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Running the Application

```bash
# Run directly
go run main.go

# Or use Make
make run
```

The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

### Testing

```bash
# Run all tests
go test -v ./...

# Or use Make
make test

# View coverage report
make coverage
```

### Available Endpoints

- `GET /health` - Health check endpoint that returns service status

Example response:
```json
{
  "status": "healthy",
  "timestamp": "2024-11-05T15:00:00Z",
  "service": "ficcc-backend"
}
```

- `GET /welcome` - Welcome message endpoint

Example response:
```json
{
  "message": "Welcome to FICCC Backend API!",
  "service": "ficcc-backend"
}
```

### Development

```bash
# Format code
make fmt

# Run linter
make lint

# Build binary
make build

# Clean build artifacts
make clean
```

## CI/CD

This project includes automated PR reviews using Claude Code. See `.github/workflows/claude-pr-review.yml` for configuration.