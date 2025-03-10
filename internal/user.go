package internal

type User struct {
	ID			int64
	FirstName 	string
	LastName	string
	Email		string
	Password	string
	Contacts	[]Contact
	Address		Address
}
