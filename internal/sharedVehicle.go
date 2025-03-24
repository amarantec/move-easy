package internal

import "time"

type SharedVehicle struct {
	ID          int64
	UserID      int64
	Latitude    float64
	Longitude   float64
	VehicleType VehicleType
	ReportedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
