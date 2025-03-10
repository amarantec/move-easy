package internal

import "time"

type UserLocation struct {
	UserID		int64
	LineID		int64
	Latitude	float64
	Longitude	float64
	TimeStamp	time.Time
}
