package sheets

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestDecodeEventsCSVParsesRows(t *testing.T) {
	events, err := decodeEventsCSV(strings.NewReader(
		"id,name,location,image_url,starts_at,ends_at,description\n" +
			"event-1,Sunday Service,Main Auditorium,https://example.com/event.jpg,12/25/2026 13:00:00,12/25/2026 15:00:00,Weekly worship gathering\n",
	))
	if err != nil {
		t.Fatalf("decodeEventsCSV returned error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.ID != "event-1" {
		t.Fatalf("expected ID event-1, got %q", event.ID)
	}
	if event.Name != "Sunday Service" {
		t.Fatalf("expected name Sunday Service, got %q", event.Name)
	}
	if event.Description != "Weekly worship gathering" {
		t.Fatalf("expected description Weekly worship gathering, got %q", event.Description)
	}
	if event.Location != "Main Auditorium" {
		t.Fatalf("expected location Main Auditorium, got %q", event.Location)
	}
	if event.ImageURL == nil || *event.ImageURL != "https://example.com/event.jpg" {
		t.Fatalf("expected image URL to be set, got %#v", event.ImageURL)
	}
	if got := event.StartsAt.Format("01/02/2006 15:04:05"); got != "12/25/2026 13:00:00" {
		t.Fatalf("unexpected StartsAt value %q", got)
	}
	if event.EndsAt == nil {
		t.Fatal("expected EndsAt to be set")
	}
	if got := event.EndsAt.Format("01/02/2006 15:04:05"); got != "12/25/2026 15:00:00" {
		t.Fatalf("unexpected EndsAt value %q", got)
	}
}

func TestDecodeEventsCSVParsesAMPMTimes(t *testing.T) {
	events, err := decodeEventsCSV(strings.NewReader(
		"id,name,location,image_url,starts_at,ends_at,description\n" +
			"event-2,Christmas Service,Sanctuary,,12/25/2026 1:00 PM,12/25/2026 4:00 PM,Christmas gathering\n",
	))
	if err != nil {
		t.Fatalf("decodeEventsCSV returned error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if got := event.StartsAt.Format("01/02/2006 15:04:05"); got != "12/25/2026 13:00:00" {
		t.Fatalf("unexpected StartsAt value %q", got)
	}
	if event.EndsAt == nil {
		t.Fatal("expected EndsAt to be set")
	}
	if got := event.EndsAt.Format("01/02/2006 15:04:05"); got != "12/25/2026 16:00:00" {
		t.Fatalf("unexpected EndsAt value %q", got)
	}
}

func TestDecodeEventsCSVParsesSingleDigitMonthDay(t *testing.T) {
	events, err := decodeEventsCSV(strings.NewReader(
		"id,name,location,image_url,starts_at,ends_at,description\n" +
			"event-3,Fall Gathering,Hall,,9/12/2026 13:00:00,9/12/2026 15:30:00,Fall event\n",
	))
	if err != nil {
		t.Fatalf("decodeEventsCSV returned error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if got := event.StartsAt.Format("1/2/2006 15:04:05"); got != "9/12/2026 13:00:00" {
		t.Fatalf("unexpected StartsAt value %q", got)
	}
	if event.EndsAt == nil {
		t.Fatal("expected EndsAt to be set")
	}
	if got := event.EndsAt.Format("1/2/2006 15:04:05"); got != "9/12/2026 15:30:00" {
		t.Fatalf("unexpected EndsAt value %q", got)
	}
}

func TestGetEventsRequiresSheetURL(t *testing.T) {
	originalValue, hadOriginal := os.LookupEnv(eventsSheetURLKey)
	if hadOriginal {
		t.Cleanup(func() {
			if err := os.Setenv(eventsSheetURLKey, originalValue); err != nil {
				t.Fatalf("restore env: %v", err)
			}
		})
	} else {
		t.Cleanup(func() {
			if err := os.Unsetenv(eventsSheetURLKey); err != nil {
				t.Fatalf("unset env: %v", err)
			}
		})
	}

	if err := os.Unsetenv(eventsSheetURLKey); err != nil {
		t.Fatalf("unset env: %v", err)
	}

	_, err := GetEvents(context.Background())
	if err == nil {
		t.Fatal("expected error when EVENTS_SHEET_CSV_URL is not set")
	}
}
