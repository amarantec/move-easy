package internal

import "time"

type UserVotes struct {
	ID		int64
	UserID	int64
	FeedbackID	int64
	VoteType	VoteType
	VotedAt		time.Time
}
