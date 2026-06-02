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

const announcementsSheetURLKey = "ANNOUNCEMENTS_SHEET_CSV_URL"

// GetAnnouncements fetches announcement records from a published Google Sheet
// CSV feed and returns only those currently within their display window.
func GetAnnouncements(ctx context.Context) ([]models.Announcement, error) {
	sheetURL := strings.TrimSpace(os.Getenv(announcementsSheetURLKey))
	if sheetURL == "" {
		return nil, fmt.Errorf("%s is not set", announcementsSheetURLKey)
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

	announcements, err := decodeAnnouncementsCSV(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse Google Sheets CSV: %w", err)
	}

	return filterActiveAnnouncements(announcements, time.Now()), nil
}

func decodeAnnouncementsCSV(reader io.Reader) ([]models.Announcement, error) {
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []models.Announcement{}, nil
	}

	header := normalizeHeader(records[0])
	announcements := make([]models.Announcement, 0, max(len(records)-1, 0))

	for rowIndex, row := range records[1:] {
		if rowIsEmpty(row) {
			continue
		}

		// rowIndex is 0-based over the data rows; the spreadsheet row number is
		// rowIndex+2 (header occupies row 1). We use it as the announcement ID.
		sheetRow := int64(rowIndex + 2)
		announcement, err := parseAnnouncementRow(header, row, sheetRow)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", sheetRow, err)
		}
		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func parseAnnouncementRow(header, row []string, id int64) (models.Announcement, error) {
	rowMap := make(map[string]string, len(header))
	for i, column := range header {
		if i < len(row) {
			rowMap[column] = strings.TrimSpace(row[i])
			continue
		}
		rowMap[column] = ""
	}

	title := rowMap["title"]
	description := rowMap["description"]
	stopRaw := rowMap["stop_displaying_at"]

	if title == "" {
		return models.Announcement{}, fmt.Errorf("missing required field title")
	}
	if description == "" {
		return models.Announcement{}, fmt.Errorf("missing required field description")
	}
	if stopRaw == "" {
		return models.Announcement{}, fmt.Errorf("missing required field stop_displaying_at")
	}

	stopAt, err := parseTime(stopRaw)
	if err != nil {
		return models.Announcement{}, fmt.Errorf("invalid stop_displaying_at %q: %w", stopRaw, err)
	}

	announcement := models.Announcement{
		ID:               id,
		Title:            title,
		Description:      description,
		StopDisplayingAt: stopAt,
	}

	if flyerURL := rowMap["flyer_url"]; flyerURL != "" {
		announcement.FlyerURL = stringPtr(flyerURL)
	}

	return announcement, nil
}

// filterActiveAnnouncements keeps announcements that have not yet reached their
// stop_displaying_at time.
func filterActiveAnnouncements(announcements []models.Announcement, now time.Time) []models.Announcement {
	active := make([]models.Announcement, 0, len(announcements))
	for _, a := range announcements {
		if !now.Before(a.StopDisplayingAt) {
			continue
		}
		active = append(active, a)
	}
	return active
}
