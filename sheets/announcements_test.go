package sheets

import (
	"strings"
	"testing"
	"time"

	"ficcc-backend/models"
)

func TestDecodeAnnouncementsCSVParsesRows(t *testing.T) {
	announcements, err := decodeAnnouncementsCSV(strings.NewReader(
		"title,description,flyer_url,stop_displaying_at\n" +
			"Easter Service,Join us for Easter,https://example.com/flyer.jpg,2026-04-30\n",
	))
	if err != nil {
		t.Fatalf("decodeAnnouncementsCSV returned error: %v", err)
	}

	if len(announcements) != 1 {
		t.Fatalf("expected 1 announcement, got %d", len(announcements))
	}

	a := announcements[0]
	if a.ID != 2 {
		t.Fatalf("expected ID 2 (spreadsheet row), got %d", a.ID)
	}
	if a.Title != "Easter Service" {
		t.Fatalf("expected title Easter Service, got %q", a.Title)
	}
	if a.Description != "Join us for Easter" {
		t.Fatalf("expected description Join us for Easter, got %q", a.Description)
	}
	if a.FlyerURL == nil || *a.FlyerURL != "https://example.com/flyer.jpg" {
		t.Fatalf("expected flyer URL to be set, got %#v", a.FlyerURL)
	}
	if got := a.StopDisplayingAt.Format("2006-01-02"); got != "2026-04-30" {
		t.Fatalf("unexpected StopDisplayingAt value %q", got)
	}
}

func TestDecodeAnnouncementsCSVOptionalFields(t *testing.T) {
	announcements, err := decodeAnnouncementsCSV(strings.NewReader(
		"title,description,flyer_url,stop_displaying_at\n" +
			"No Flyer,Body text,,2026-04-30\n",
	))
	if err != nil {
		t.Fatalf("decodeAnnouncementsCSV returned error: %v", err)
	}

	a := announcements[0]
	if a.FlyerURL != nil {
		t.Fatalf("expected nil FlyerURL, got %#v", a.FlyerURL)
	}
}

func TestDecodeAnnouncementsCSVMissingRequiredField(t *testing.T) {
	_, err := decodeAnnouncementsCSV(strings.NewReader(
		"title,description,flyer_url,stop_displaying_at\n" +
			"Has Title,Has body,,\n",
	))
	if err == nil {
		t.Fatal("expected error for missing stop_displaying_at, got nil")
	}
}

func TestFilterActiveAnnouncements(t *testing.T) {
	now := time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	announcements := []struct {
		name   string
		stopAt time.Time
		want   bool
	}{
		{"still active", future, true},
		{"already stopped", past, false},
		{"stops exactly now", now, false},
	}

	for _, tc := range announcements {
		in := []models.Announcement{{
			Title:            tc.name,
			StopDisplayingAt: tc.stopAt,
		}}
		got := filterActiveAnnouncements(in, now)
		if (len(got) == 1) != tc.want {
			t.Fatalf("%s: expected active=%v, got %d results", tc.name, tc.want, len(got))
		}
	}
}
