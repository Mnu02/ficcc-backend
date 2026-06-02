package sheets

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"ficcc-backend/models"
)

const (
	eventsSheetURLKey = "EVENTS_SHEET_CSV_URL"
	httpTimeout       = 10 * time.Second
)

// GetEvents fetches event records from a published Google Sheet CSV feed.
func GetEvents(ctx context.Context) ([]models.Event, error) {
	sheetURL := strings.TrimSpace(os.Getenv(eventsSheetURLKey))
	if sheetURL == "" {
		return nil, fmt.Errorf("%s is not set", eventsSheetURLKey)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sheetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create Google Sheets request: %w", err)
	}

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Google Sheets CSV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch Google Sheets CSV: unexpected status %d", resp.StatusCode)
	}

	events, err := decodeEventsCSV(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse Google Sheets CSV: %w", err)
	}

	return events, nil
}

func decodeEventsCSV(reader io.Reader) ([]models.Event, error) {
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []models.Event{}, nil
	}

	header := normalizeHeader(records[0])
	events := make([]models.Event, 0, max(len(records)-1, 0))

	for rowIndex, row := range records[1:] {
		if rowIsEmpty(row) {
			continue
		}

		// rowIndex is 0-based over the data rows; the spreadsheet row number is
		// rowIndex+2 (header occupies row 1). We use it as the event ID.
		sheetRow := rowIndex + 2
		event, err := parseEventRow(header, row, sheetRow)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", sheetRow, err)
		}
		events = append(events, event)
	}

	return events, nil
}

func parseEventRow(header, row []string, id int) (models.Event, error) {
	rowMap := make(map[string]string, len(header))
	for i, column := range header {
		if i < len(row) {
			rowMap[column] = strings.TrimSpace(row[i])
			continue
		}
		rowMap[column] = ""
	}

	name := rowMap["name"]
	location := rowMap["location"]
	startsAtRaw := rowMap["starts_at"]
	description := rowMap["description"]

	if name == "" {
		return models.Event{}, fmt.Errorf("missing required field name")
	}
	if location == "" {
		return models.Event{}, fmt.Errorf("missing required field location")
	}
	if startsAtRaw == "" {
		return models.Event{}, fmt.Errorf("missing required field starts_at")
	}
	if description == "" {
		return models.Event{}, fmt.Errorf("missing required field description")
	}

	startsAt, err := parseTime(startsAtRaw)
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid starts_at %q: %w", startsAtRaw, err)
	}

	event := models.Event{
		ID:          int64(id),
		Name:        name,
		Description: description,
		Location:    location,
		StartsAt:    startsAt,
	}

	if imageURL := rowMap["image_url"]; imageURL != "" {
		event.ImageURL = stringPtr(imageURL)
	}

	if endsAtRaw := rowMap["ends_at"]; endsAtRaw != "" {
		endsAt, err := parseTime(endsAtRaw)
		if err != nil {
			return models.Event{}, fmt.Errorf("invalid ends_at %q: %w", endsAtRaw, err)
		}
		event.EndsAt = &endsAt
	}

	return event, nil
}

func parseTime(value string) (time.Time, error) {
	layouts := []string{
		"1/2/2006 3:04 PM",
		"1/2/2006 3:04:05 PM",
		"1/2/2006 15:04:05",
		"01/02/2006 3:04 PM",
		"01/02/2006 3:04:05 PM",
		"01/02/2006 15:04:05",
		"1/2/2006",
		"01/02/2006",
		time.RFC3339,
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC1123Z,
		time.RFC1123,
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("supported formats: 1/2/2006 3:04 PM, 1/2/2006 3:04:05 PM, 1/2/2006 15:04:05, 01/02/2006 3:04 PM, 01/02/2006 3:04:05 PM, 01/02/2006 15:04:05, 1/2/2006, 01/02/2006, RFC3339, 2006-01-02 15:04, 2006-01-02")
}

func normalizeHeader(row []string) []string {
	normalized := make([]string, len(row))
	for i, value := range row {
		normalized[i] = strings.ToLower(strings.TrimSpace(value))
	}
	return normalized
}

func rowIsEmpty(row []string) bool {
	for _, value := range row {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func stringPtr(value string) *string {
	return &value
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
