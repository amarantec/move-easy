package internal

import "time"

type Occurrence struct {
	ID				int64
	UserID			int64
	Type			OccurrenceType
	Description		string
	TimeStamp 		time.Time
	Confirmation	int64
}
