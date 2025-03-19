package internal

import "time"

type Address struct {
	ID		        int64
	UserID	        int64
	Street 	        string
	Number	        string
	CEP		        string
	Neighborhood    string
	City	        string
	State	        string
    CreatedAt       time.Time
    UpdatedAt       *time.Time
    DeletedAt       *time.Time
}
