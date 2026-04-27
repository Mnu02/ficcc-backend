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
- `GET /events` - gets all events from Google Sheets
- `GET /videos` - gets recent public videos from the configured YouTube channel

### Events via Google Sheets

Events can be served from a Google Sheet instead of Supabase. This is intended for non-technical staff who are comfortable editing spreadsheet rows.

1. Create a Google Sheet with a tab for events.
2. Use this exact header row:

```text
id,name,location,image_url,starts_at,ends_at,description
```

3. Publish the sheet as CSV or copy a public CSV export URL for the events tab.
4. Add this to your `.env`:

```bash
EVENTS_SHEET_CSV_URL=https://docs.google.com/spreadsheets/d/.../pub?output=csv
```

Supported date/time formats for `starts_at` and `ends_at`:
- `12/25/2026 4:00 PM`
- `12/25/2026 4:00:00 PM`
- `12/25/2026 13:00:00`
- `2026-05-01T18:30:00Z`
- `2026-05-01 18:30`
- `2026-05-01`

Notes:
- `id`, `name`, `location`, `description` and `starts_at` are required.
- `image_url` and `ends_at` are optional.
- If `DATABASE_URL` is not set, the server still starts and `/events` will work, but database-backed endpoints like `/sermons` will not.

### Videos via YouTube

Videos are served from the YouTube Data API v3.

1. Create a Google Cloud API key with YouTube Data API v3 enabled.
2. Add this to your `.env`:

```bash
YOUTUBE_API_KEY=your-api-key
YOUTUBE_CHANNEL_ID=your-channel-id
```

You can also set `YOUTUBE_UPLOADS_PLAYLIST_ID` directly. If it is set, the backend skips the channel lookup and reads videos from that uploads playlist.

Optional query param:

```text
GET /videos?maxResults=10
```

`maxResults` defaults to 25 and is capped at 50.

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
