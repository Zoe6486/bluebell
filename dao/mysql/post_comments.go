package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"bluebell/models"
)

// PostCommentStore defines the data access interface for post comments.
type PostCommentStore interface {
	Create(ctx context.Context, c *models.PostComment) error
	GetByID(ctx context.Context, commentID uint64) (*models.PostComment, error)
	ListByPost(ctx context.Context, params *models.CommentListParams) ([]*models.CommentDetail, int64, error)
	Delete(ctx context.Context, commentID uint64, userID uint64) error
	IncrLikeCount(ctx context.Context, commentID uint64, delta int) error
}

type postCommentDao struct {
	db *sqlx.DB
}

// NewPostCommentDao constructs a PostCommentStore backed by MySQL.
func NewPostCommentDao(db *sqlx.DB) PostCommentStore {
	return &postCommentDao{db: db}
}

func (d *postCommentDao) Create(ctx context.Context, c *models.PostComment) error {
	const query = `
		INSERT INTO post_comments (post_id, user_id, content, status)
		VALUES (:post_id, :user_id, :content, :status)`

	res, err := d.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = uint64(id)
	return nil
}

func (d *postCommentDao) GetByID(ctx context.Context, commentID uint64) (*models.PostComment, error) {
	const query = `
		SELECT id, post_id, user_id, content, like_count, status, create_time
		FROM post_comments
		WHERE id = ? AND status = ?`

	var c models.PostComment
	if err := d.db.GetContext(ctx, &c, query, commentID, models.CommentStatusNormal); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ListByPost returns a paginated, time-ordered list of comments for a post,
// joined with author info.
func (d *postCommentDao) ListByPost(ctx context.Context, params *models.CommentListParams) ([]*models.CommentDetail, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}
	// 如果Page和PageSize没有默认值，那model里就该加 binding:"required"
	offset := (params.Page - 1) * params.PageSize

	var total int64
	const countSQL = `
		SELECT COUNT(*) FROM post_comments
		WHERE post_id = ? AND status = ?`
	if err := d.db.GetContext(ctx, &total, countSQL, params.PostID, models.CommentStatusNormal); err != nil {
		return nil, 0, err
	}

	const dataSQL = `
		SELECT
			c.id, c.post_id, c.user_id, c.content, c.like_count, c.status, c.create_time,
			u.username    AS author_name
		FROM post_comments c
		LEFT JOIN user u ON u.user_id = c.user_id
		WHERE c.post_id = ? AND c.status = ?
		ORDER BY c.create_time DESC
		LIMIT ? OFFSET ?`

	var comments []*models.CommentDetail
	if err := d.db.SelectContext(ctx, &comments, dataSQL,
		params.PostID, models.CommentStatusNormal, params.PageSize, offset,
	); err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// Delete soft-deletes a comment. userID is checked so only the owner can delete.
func (d *postCommentDao) Delete(ctx context.Context, commentID uint64, userID uint64) error {
	const query = `
		UPDATE post_comments
		SET status = ?
		WHERE id = ? AND user_id = ? AND status = ?`

	res, err := d.db.ExecContext(ctx, query,
		models.CommentStatusDeleted, commentID, userID, models.CommentStatusNormal)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// IncrLikeCount atomically increments (or decrements with delta=-1) like_count.
func (d *postCommentDao) IncrLikeCount(ctx context.Context, commentID uint64, delta int) error {
	const query = `
		UPDATE post_comments
		SET like_count = like_count + ?
		WHERE id = ? AND status = ?`

	_, err := d.db.ExecContext(ctx, query, delta, commentID, models.CommentStatusNormal)
	return err
}
