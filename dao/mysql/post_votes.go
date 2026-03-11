package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"bluebell/models"
)

// PostVoteStore defines the data access interface for post votes.
type PostVoteStore interface {
	// Upsert inserts or updates a vote. vote_type=0 means cancel (delete the row).
	Upsert(ctx context.Context, params *models.VotePostParams) error
	// GetUserVote returns the current vote a user has on a post, or nil if none.
	GetUserVote(ctx context.Context, postID, userID int64) (*models.PostVote, error)
	// GetSummary returns aggregated like/dislike counts for a post.
	GetSummary(ctx context.Context, postID uint64) (*models.PostVoteSummary, error)
	// GetSummaryBatch returns vote summaries for multiple posts in one query.
	GetSummaryBatch(ctx context.Context, postIDs []uint64) (map[uint64]*models.PostVoteSummary, error)
}

type postVoteDao struct {
	db *sqlx.DB
}

// NewPostVoteDao constructs a PostVoteStore backed by MySQL.
func NewPostVoteDao(db *sqlx.DB) PostVoteStore {
	return &postVoteDao{db: db}
}

func (d *postVoteDao) Upsert(ctx context.Context, params *models.VotePostParams) error {
	if params.VoteType == models.VoteTypeCancel {
		// Cancel vote — delete the row entirely
		const query = `DELETE FROM post_votes WHERE post_id = ? AND user_id = ?`
		_, err := d.db.ExecContext(ctx, query, params.PostID, params.UserID)
		return err
	}

	// INSERT ... ON DUPLICATE KEY UPDATE handles both new votes and changed votes
	const query = `
		INSERT INTO post_votes (post_id, user_id, vote_type)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE vote_type = VALUES(vote_type), update_time = CURRENT_TIMESTAMP`

	_, err := d.db.ExecContext(ctx, query, params.PostID, params.UserID, params.VoteType)
	return err
}

func (d *postVoteDao) GetUserVote(ctx context.Context, postID, userID int64) (*models.PostVote, error) {
	const query = `
		SELECT id, post_id, user_id, vote_type, create_time, update_time
		FROM post_votes
		WHERE post_id = ? AND user_id = ?`

	var v models.PostVote
	if err := d.db.GetContext(ctx, &v, query, postID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // no vote is a valid state, not an error
		}
		return nil, err
	}
	return &v, nil
}

func (d *postVoteDao) GetSummary(ctx context.Context, postID uint64) (*models.PostVoteSummary, error) {
	const query = `
		SELECT
			post_id,
			SUM(CASE WHEN vote_type =  1 THEN 1 ELSE 0 END) AS like_count,
			SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END) AS dislike_count
		FROM post_votes
		WHERE post_id = ?
		GROUP BY post_id`

	var summary models.PostVoteSummary
	if err := d.db.GetContext(ctx, &summary, query, postID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Post exists but has zero votes — return zeroed struct
			return &models.PostVoteSummary{PostID: postID}, nil
		}
		return nil, err
	}
	return &summary, nil
}

func (d *postVoteDao) GetSummaryBatch(ctx context.Context, postIDs []uint64) (map[uint64]*models.PostVoteSummary, error) {
	if len(postIDs) == 0 {
		return map[uint64]*models.PostVoteSummary{}, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			post_id,
			SUM(CASE WHEN vote_type =  1 THEN 1 ELSE 0 END) AS like_count,
			SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END) AS dislike_count
		FROM post_votes
		WHERE post_id IN (?)
		GROUP BY post_id`, postIDs)
	if err != nil {
		return nil, err
	}
	query = d.db.Rebind(query)

	var rows []models.PostVoteSummary
	if err = d.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make(map[uint64]*models.PostVoteSummary, len(rows))
	for i := range rows {
		r := rows[i]
		result[r.PostID] = &r
	}
	return result, nil
}
