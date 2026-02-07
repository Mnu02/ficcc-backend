package models

import "time"

// PrayerRequest represents a prayer request submitted by a user
type PrayerRequest struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      string    `json:"user_id"`
	Body        string    `json:"body"`
	IsAnonymous bool      `json:"is_anonymous"`
	Status      string    `json:"status"`
}
