package internal

type BusLine struct {
	ID			int64
	Name		string
	Schedules 	[]BusSchedules
}
