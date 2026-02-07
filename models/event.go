package models

import "time"

// Event represents a church event
type Event struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Location string     `json:"location"`
	ImageURL *string    `json:"image_url"`
	StartsAt time.Time  `json:"starts_at"`
	EndsAt   *time.Time `json:"ends_at"`
}
