package models

import "time"

// PostVote represents a user's vote on a post
type PostVote struct {
	ID        uint64    `db:"id"`
	PostID    uint64    `db:"post_id"`
	UserID    uint64    `db:"user_id"`
	VoteType  int8      `db:"vote_type"`
	CreatedAt time.Time `db:"create_time"`
	UpdatedAt time.Time `db:"update_time"`
}

// VoteType constants
const (
	VoteTypeLike    int8 = 1
	VoteTypeDislike int8 = -1
	VoteTypeCancel  int8 = 0
)

// VotePostParams is used when a user submits a vote
type VotePostParams struct {
	PostID   int64 `json:"post_id" binding:"required"`
	UserID   int64 `json:"user_id"`
	VoteType int8  `json:"vote_type" binding:"oneof=-1 0 1"`
}

// PostVoteSummary holds aggregated vote counts for a post
type PostVoteSummary struct {
	PostID       uint64 `db:"post_id"`
	LikeCount    int64  `db:"like_count"`
	DislikeCount int64  `db:"dislike_count"`
}
