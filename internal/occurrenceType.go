package internal

type OccurrenceType int

const (
	STOPPED_TRAFFIC OccurrenceType = iota
	ACCIDENT
	LOCKED_BUS
	ITINERARY_CHANGE
	OTHER
)
	
