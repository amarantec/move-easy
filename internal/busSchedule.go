package internal

import "time"

type BusSchedules struct {
	ID        int64
	BusLineID int64
	DayOfWeek string
	StartTime *time.Time
	EndTime   *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
