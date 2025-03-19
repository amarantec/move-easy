package internal

import "time"

type Contact struct {
	ID			int64
	UserID		int64
	Name		string
	DDI			string
	DDD			string
	PhoneNumber	string
    CreatedAt   time.Time
    Updated     *time.Time
    DeletedAt   *time.Time
}
