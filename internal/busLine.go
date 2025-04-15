package internal

import "time"

type BusLine struct {
	ID        int64
	Name      string
	BusInit   BusStop
	BusEnd    BusStop
	Schedules []BusSchedules
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
