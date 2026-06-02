package models

import "time"

// Announcement represents a church announcement sourced from a Google Sheet.
//
// ID is not entered by the (non-technical) sheet editor; it is the underlying
// spreadsheet row number, assigned when the sheet is read, so the frontend has
// a stable key for rendering.
type Announcement struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	FlyerURL         *string   `json:"flyer_url"` // optional
	StopDisplayingAt time.Time `json:"stop_displaying_at"`
}
