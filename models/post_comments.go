package models

import "time"

// PostComment represents a top-level comment on a post
type PostComment struct {
	ID         uint64    `db:"id"`
	PostID     uint64    `db:"post_id"`
	UserID     uint64    `db:"user_id"`
	Content    string    `db:"content"`
	LikeCount  uint32    `db:"like_count"`
	Status     int8      `db:"status"`
	CreateTime time.Time `db:"create_time"`
}

// CommentStatus constants
const (
	CommentStatusNormal  int8 = 1
	CommentStatusDeleted int8 = 2
)

// CommentDetail is used for API responses, with author info joined in
type CommentDetail struct {
	*PostComment
	AuthorName string `db:"author_name"`
	// AuthorAvatar string `db:"author_avatar"`
}

// CreateCommentParams holds data needed to create a comment
type CreateCommentParams struct {
	PostID  uint64 `json:"post_id"`
	UserID  uint64 `json:"user_id"`
	Content string `json:"content" binding:"required,max=2048"` // body 参数（content）→ 在 body 里，必须用 binding 验证
}

// CommentListParams holds pagination params for listing comments
type CommentListParams struct {
	// PostID   uint64 `json:"post_id" form:"post_id" binding:"required"`删掉required
	// CommentListParams 里加它就不合适，因为 post_id 根本不从 query 参数来，Gin 绑定时永远找不到它，就永远报 400。
	PostID   uint64 `json:"post_id" form:"post_id"`
	Page     int64  `json:"page" form:"page"`           // query 参数, 因为page, pageSize在dao/mysql/post_comments代码里有默认值兜底值
	PageSize int64  `json:"page_size" form:"page_size"` // query 参数, 如果没有默认值，那就该加 binding:"required" 了，不传就报 400 让用户必须传。
}
