package models

import "time"

type Video struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PublishedAt  time.Time `json:"published_at"`
	ThumbnailURL string    `json:"thumbnail_url"`
	WatchURL     string    `json:"watch_url"`
	EmbedURL     string    `json:"embed_url"`
}
