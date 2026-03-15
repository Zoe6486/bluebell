package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"bluebell/logic"
	"bluebell/models"
)

type postDao struct {
	db *sqlx.DB
}

func NewPostDao(db *sqlx.DB) logic.PostStore {
	return &postDao{db: db}
}

func (d *postDao) Create(ctx context.Context, p *models.Post) error {
	const query = `
		INSERT INTO post (post_id, title, content, author_id, community_id, status)
		VALUES (:post_id, :title, :content, :author_id, :community_id, :status)`
	_, err := d.db.NamedExecContext(ctx, query, p)
	return err
}

func (d *postDao) GetByID(ctx context.Context, postID int64) (*models.Post, error) {
	const query = `
		SELECT id, post_id, title, content, author_id, community_id, status, create_time, update_time
		FROM post
		WHERE post_id = ? AND status != ?`

	var p models.Post
	if err := d.db.GetContext(ctx, &p, query, postID, models.PostStatusDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logic.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (d *postDao) GetDetail(ctx context.Context, postID int64) (*models.PostDetail, error) {
	const query = `
		SELECT
			p.id, p.post_id, p.title, p.content, p.author_id, p.community_id,
			p.status, p.create_time, p.update_time,
			u.username                                                          AS author_name,
			c.community_name                                                    AS community_name,
			COALESCE(SUM(CASE WHEN pv.vote_type =  1 THEN 1 ELSE 0 END), 0)   AS like_count,
			COALESCE(SUM(CASE WHEN pv.vote_type = -1 THEN 1 ELSE 0 END), 0)   AS dislike_count
		FROM post p
		LEFT JOIN user       u  ON u.user_id      = p.author_id
		LEFT JOIN community  c  ON c.community_id = p.community_id
		LEFT JOIN post_votes pv ON pv.post_id      = p.post_id
		WHERE p.post_id = ? AND p.status != ?
		GROUP BY p.id`

	var detail models.PostDetail
	detail.Post = &models.Post{}
	if err := d.db.GetContext(ctx, &detail, query, postID, models.PostStatusDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logic.ErrNotFound
		}
		return nil, err
	}
	return &detail, nil
}

func (d *postDao) List(ctx context.Context, params *models.PostListParams) ([]*models.Post, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}
	offset := (params.Page - 1) * params.PageSize

	orderBy := "p.create_time DESC"
	if params.OrderBy == "score" {
		orderBy = "like_count DESC, p.create_time DESC"
	}

	whereClause := "WHERE p.status = ?"
	args := []any{models.PostStatusNormal}

	if params.CommunityID > 0 {
		whereClause += " AND p.community_id = ?"
		args = append(args, params.CommunityID)
	}

	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM post p %s`, whereClause)
	var total int64
	if err := d.db.GetContext(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	dataSQL := fmt.Sprintf(`
		SELECT
			p.id, p.post_id, p.title, p.content, p.author_id, p.community_id,
			p.status, p.create_time, p.update_time,
			COALESCE(SUM(CASE WHEN pv.vote_type = 1 THEN 1 ELSE 0 END), 0) AS like_count
		FROM post p
		LEFT JOIN post_votes pv ON pv.post_id = p.post_id
		%s
		GROUP BY p.id
		ORDER BY %s
		LIMIT ? OFFSET ?`, whereClause, orderBy)

	args = append(args, params.PageSize, offset)
	var posts []*models.Post
	if err := d.db.SelectContext(ctx, &posts, dataSQL, args...); err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (d *postDao) Delete(ctx context.Context, postID int64) error {
	const query = `UPDATE post SET status = ? WHERE post_id = ?`
	res, err := d.db.ExecContext(ctx, query, models.PostStatusDeleted, postID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return logic.ErrNotFound
	}
	return nil
}
