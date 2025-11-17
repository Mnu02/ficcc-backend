package models

import "time"

// Sermon represents a sermon record from the database
type Sermon struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Preacher     string    `json:"preacher"`
	ScriptureRef string    `json:"scripture_ref"`
	SermonDate   time.Time `json:"sermon_date"`
	SermonSeries *string   `json:"sermon_series"` // Nullable field
	YouTubeLink  string    `json:"youtube_link"`
}
