package models

import "time"

// Announcement represents a church announcement
type Announcement struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	BodyHTML  string    `json:"body_html"`
	ImageURL  *string   `json:"image_url"`
}
