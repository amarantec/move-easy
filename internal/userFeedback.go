package internal

import "time"

type UserFeedback struct {
	ID			int64
	UserID		int64
	Latitude	float64
	Longitude	float64
	Description	string
	Category	string
	CreatedAt	time.Time
}
