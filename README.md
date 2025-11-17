# ficcc-backend

A Golang backend service for FICCC.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Make (optional, for using Makefile commands)
- Supabase account with a project set up

### Running the Application

The application will automatically load environment variables from the `.env` file:

```bash
# Run directly
go run main.go

# Or use Make
make run
```

The server will start on port 8080 by default (or whatever you set in your `.env` file).

**Important:** Make sure your `.env` file exists with a valid `DATABASE_URL` before running the application, otherwise it will fail to start.

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

- `GET /sermons` - gets all sermons

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