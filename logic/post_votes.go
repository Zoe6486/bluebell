package logic

import (
	"context"
	"time"

	"bluebell/dao/mysql"
	"bluebell/models"
)

const voteExpiryDuration = 7 * 24 * time.Hour // posts expire for voting after 1 week

// PostVoteLogic handles business rules for voting.
type PostVoteLogic struct {
	voteDao mysql.PostVoteStore
	postDao mysql.PostStore
}

// NewPostVoteLogic wires up vote logic with its dependencies.
func NewPostVoteLogic(voteDao mysql.PostVoteStore, postDao mysql.PostStore) *PostVoteLogic {
	return &PostVoteLogic{
		voteDao: voteDao,
		postDao: postDao,
	}
}

// VotePost applies a user's vote to a post, enforcing expiry and idempotency.
func (l *PostVoteLogic) VotePost(ctx context.Context, params *models.VotePostParams) error {
	// 1. Check post exists and is not expired
	post, err := l.postDao.GetByID(ctx, params.PostID)
	if err != nil {
		return err
	}
	if time.Since(post.CreateTime) > voteExpiryDuration {
		return ErrPostTooOld
	}

	// 2. Check existing vote to detect no-op (avoid unnecessary DB write)
	existing, err := l.voteDao.GetUserVote(ctx, params.PostID, params.UserID)
	if err != nil {
		return err
	}

	// If the user is casting the same vote again, treat it as idempotent success
	if existing != nil && existing.VoteType == params.VoteType {
		return nil
	}

	// 3. Upsert (handles new vote, changed vote, and cancel)
	return l.voteDao.Upsert(ctx, params)
}

// GetPostVoteSummary returns like/dislike counts for a single post.
func (l *PostVoteLogic) GetPostVoteSummary(ctx context.Context, postID uint64) (*models.PostVoteSummary, error) {
	return l.voteDao.GetSummary(ctx, postID)
}

// GetUserVoteStatus returns what vote (if any) the user has on a post.
// Returns nil, nil if the user has not voted.
func (l *PostVoteLogic) GetUserVoteStatus(ctx context.Context, postID, userID int64) (*models.PostVote, error) {
	return l.voteDao.GetUserVote(ctx, postID, userID)
}
