package youtube

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUploadsPlaylistID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/channels" {
			t.Fatalf("expected /channels path, got %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("part"); got != "contentDetails" {
			t.Fatalf("expected part contentDetails, got %q", got)
		}
		if got := r.URL.Query().Get("id"); got != "channel-1" {
			t.Fatalf("expected channel ID channel-1, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"items": [{
				"contentDetails": {
					"relatedPlaylists": {
						"uploads": "uploads-1"
					}
				}
			}]
		}`))
	}))
	defer server.Close()

	playlistID, err := getUploadsPlaylistID(context.Background(), server.Client(), server.URL, "test-key", "channel-1")
	if err != nil {
		t.Fatalf("getUploadsPlaylistID returned error: %v", err)
	}
	if playlistID != "uploads-1" {
		t.Fatalf("expected uploads-1, got %q", playlistID)
	}
}

func TestGetPlaylistVideosMapsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/playlistItems" {
			t.Fatalf("expected /playlistItems path, got %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("playlistId"); got != "uploads-1" {
			t.Fatalf("expected playlist ID uploads-1, got %q", got)
		}
		if got := r.URL.Query().Get("maxResults"); got != "10" {
			t.Fatalf("expected maxResults 10, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"items": [{
				"snippet": {
					"publishedAt": "2026-04-01T12:00:00Z",
					"title": "Sunday Service",
					"description": "Weekly service",
					"thumbnails": {
						"default": {"url": "https://example.com/default.jpg"},
						"high": {"url": "https://example.com/high.jpg"}
					},
					"resourceId": {
						"videoId": "abc123"
					}
				},
				"contentDetails": {
					"videoId": "abc123",
					"videoPublishedAt": "2026-04-02T12:00:00Z"
				}
			}]
		}`))
	}))
	defer server.Close()

	videos, err := getPlaylistVideos(context.Background(), server.Client(), server.URL, "test-key", "uploads-1", 10)
	if err != nil {
		t.Fatalf("getPlaylistVideos returned error: %v", err)
	}
	if len(videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(videos))
	}

	video := videos[0]
	if video.ID != "abc123" {
		t.Fatalf("expected ID abc123, got %q", video.ID)
	}
	if video.Title != "Sunday Service" {
		t.Fatalf("expected title Sunday Service, got %q", video.Title)
	}
	if video.ThumbnailURL != "https://example.com/high.jpg" {
		t.Fatalf("expected high thumbnail, got %q", video.ThumbnailURL)
	}
	if video.WatchURL != "https://www.youtube.com/watch?v=abc123" {
		t.Fatalf("unexpected watch URL %q", video.WatchURL)
	}
	if video.EmbedURL != "https://www.youtube.com/embed/abc123" {
		t.Fatalf("unexpected embed URL %q", video.EmbedURL)
	}
	if got := video.PublishedAt.Format("2006-01-02T15:04:05Z"); got != "2026-04-02T12:00:00Z" {
		t.Fatalf("unexpected published date %q", got)
	}
}

func TestNormalizeMaxResults(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "default", in: 0, want: defaultMaxResults},
		{name: "negative", in: -1, want: defaultMaxResults},
		{name: "valid", in: 12, want: 12},
		{name: "cap", in: 500, want: maxAllowedResults},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeMaxResults(tt.in); got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
