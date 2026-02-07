package models

// Leader represents a church leader
type Leader struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Bio         *string  `json:"bio"`
	ImageURL    string   `json:"image_url"`
	Email       *string  `json:"email"`
	PhoneNumber *string  `json:"phone_number"`
	SortOrder   *float64 `json:"sort_order"`
}
