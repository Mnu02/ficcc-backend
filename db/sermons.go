package db

import (
	"context"
	"fmt"

	"ficcc-backend/models"
)

// GetSermons retrieves all sermons from the database and returns them as a slice
func GetSermons(ctx context.Context) ([]models.Sermon, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Get all sermons
	rows, err := DB.Query(ctx, "SELECT * FROM sermons")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Slice to hold all sermons
	var sermons []models.Sermon

	// Process rows
	for rows.Next() {
		var sermon models.Sermon

		if err := rows.Scan(
			&sermon.ID,
			&sermon.Title,
			&sermon.Preacher,
			&sermon.ScriptureRef,
			&sermon.SermonDate, // Scan date directly into time.Time
			&sermon.SermonSeries,
			&sermon.YouTubeLink,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sermon row: %w", err)
		}

		sermons = append(sermons, sermon)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return sermons, nil
}
