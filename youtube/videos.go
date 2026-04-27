package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"ficcc-backend/models"
)

const (
	apiKeyEnv            = "YOUTUBE_API_KEY"
	channelIDEnv         = "YOUTUBE_CHANNEL_ID"
	uploadsPlaylistIDEnv = "YOUTUBE_UPLOADS_PLAYLIST_ID"
	defaultMaxResults    = 25
	maxAllowedResults    = 50
	youtubeAPIBaseURL    = "https://www.googleapis.com/youtube/v3"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetVideos fetches public videos from the configured YouTube channel.
func GetVideos(ctx context.Context, maxResults int) ([]models.Video, error) {
	apiKey := strings.TrimSpace(os.Getenv(apiKeyEnv))
	if apiKey == "" {
		return nil, fmt.Errorf("%s is not set", apiKeyEnv)
	}

	playlistID := strings.TrimSpace(os.Getenv(uploadsPlaylistIDEnv))
	if playlistID == "" {
		channelID := strings.TrimSpace(os.Getenv(channelIDEnv))
		if channelID == "" {
			return nil, fmt.Errorf("%s or %s must be set", uploadsPlaylistIDEnv, channelIDEnv)
		}

		var err error
		playlistID, err = getUploadsPlaylistID(ctx, httpClient, youtubeAPIBaseURL, apiKey, channelID)
		if err != nil {
			return nil, err
		}
	}

	return getPlaylistVideos(ctx, httpClient, youtubeAPIBaseURL, apiKey, playlistID, normalizeMaxResults(maxResults))
}

func getUploadsPlaylistID(ctx context.Context, client *http.Client, baseURL, apiKey, channelID string) (string, error) {
	values := url.Values{}
	values.Set("part", "contentDetails")
	values.Set("id", channelID)
	values.Set("key", apiKey)

	var response channelListResponse
	if err := getJSON(ctx, client, baseURL+"/channels?"+values.Encode(), &response); err != nil {
		return "", fmt.Errorf("fetch YouTube channel: %w", err)
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("YouTube channel %q not found", channelID)
	}

	playlistID := strings.TrimSpace(response.Items[0].ContentDetails.RelatedPlaylists.Uploads)
	if playlistID == "" {
		return "", fmt.Errorf("YouTube channel %q has no uploads playlist", channelID)
	}

	return playlistID, nil
}

func getPlaylistVideos(ctx context.Context, client *http.Client, baseURL, apiKey, playlistID string, maxResults int) ([]models.Video, error) {
	values := url.Values{}
	values.Set("part", "snippet,contentDetails")
	values.Set("playlistId", playlistID)
	values.Set("maxResults", strconv.Itoa(maxResults))
	values.Set("key", apiKey)

	var response playlistItemsResponse
	if err := getJSON(ctx, client, baseURL+"/playlistItems?"+values.Encode(), &response); err != nil {
		return nil, fmt.Errorf("fetch YouTube playlist videos: %w", err)
	}

	videos := make([]models.Video, 0, len(response.Items))
	for _, item := range response.Items {
		videoID := strings.TrimSpace(item.ContentDetails.VideoID)
		if videoID == "" {
			videoID = strings.TrimSpace(item.Snippet.ResourceID.VideoID)
		}
		if videoID == "" {
			continue
		}

		publishedAt := item.ContentDetails.VideoPublishedAt
		if publishedAt.IsZero() {
			publishedAt = item.Snippet.PublishedAt
		}

		videos = append(videos, models.Video{
			ID:           videoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  publishedAt,
			ThumbnailURL: bestThumbnailURL(item.Snippet.Thumbnails),
			WatchURL:     "https://www.youtube.com/watch?v=" + videoID,
			EmbedURL:     "https://www.youtube.com/embed/" + videoID,
		})
	}

	return videos, nil
}

func getJSON(ctx context.Context, client *http.Client, requestURL string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func normalizeMaxResults(maxResults int) int {
	if maxResults <= 0 {
		return defaultMaxResults
	}
	if maxResults > maxAllowedResults {
		return maxAllowedResults
	}
	return maxResults
}

func bestThumbnailURL(thumbnails map[string]thumbnail) string {
	for _, key := range []string{"maxres", "standard", "high", "medium", "default"} {
		if thumbnail, ok := thumbnails[key]; ok && thumbnail.URL != "" {
			return thumbnail.URL
		}
	}
	return ""
}

type channelListResponse struct {
	Items []struct {
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type playlistItemsResponse struct {
	Items []playlistItem `json:"items"`
}

type playlistItem struct {
	Snippet struct {
		PublishedAt time.Time            `json:"publishedAt"`
		Title       string               `json:"title"`
		Description string               `json:"description"`
		Thumbnails  map[string]thumbnail `json:"thumbnails"`
		ResourceID  struct {
			VideoID string `json:"videoId"`
		} `json:"resourceId"`
	} `json:"snippet"`
	ContentDetails struct {
		VideoID          string    `json:"videoId"`
		VideoPublishedAt time.Time `json:"videoPublishedAt"`
	} `json:"contentDetails"`
}

type thumbnail struct {
	URL string `json:"url"`
}
