package internal

import "time"

type BusStop struct {
	ID        int64
	Name      string
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
