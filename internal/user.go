package internal

import "time"

type User struct {
	ID			int64
	FirstName 	string
	LastName	string
	Email		string
	Password	string
	Contacts	[]Contact
	Address		Address
    CreatedAt   time.Time
    UpdatedAt   *time.Time
    DeletedAt   *time.Time
}
