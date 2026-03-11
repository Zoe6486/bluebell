// package logic

// import (
// 	"context"
// 	"errors"
// 	"time"

// 	"github.com/bwmarrin/snowflake"

// 	"bluebell/dao/mysql"
// 	"bluebell/models"
// )

// // PostLogic handles business rules for posts.
// type PostLogic struct {
// 	postDao mysql.PostStore
// 	voteDao mysql.PostVoteStore
// 	node    *snowflake.Node // for generating post_id
// }

// // NewPostLogic wires up the logic layer.
// // Accepts interfaces so tests can inject mocks.
// func NewPostLogic(postDao mysql.PostStore, voteDao mysql.PostVoteStore, node *snowflake.Node) *PostLogic {
// 	return &PostLogic{
// 		postDao: postDao,
// 		voteDao: voteDao,
// 		node:    node,
// 	}
// }

// // CreatePost validates input, generates a snowflake ID, and persists the post.
// func (l *PostLogic) CreatePost(ctx context.Context, params *models.CreatePostParams) (*models.Post, error) {
// 	post := &models.Post{
// 		PostID:      l.node.Generate().Int64(),
// 		Title:       params.Title,
// 		Content:     params.Content,
// 		AuthorID:    params.AuthorID,
// 		CommunityID: params.CommunityID,
// 		Status:      models.PostStatusNormal,
// 	}

// 	if err := l.postDao.Create(ctx, post); err != nil {
// 		return nil, err
// 	}
// 	return post, nil
// }

// // GetPostDetail fetches the full post with author, community and vote info.
// func (l *PostLogic) GetPostDetail(ctx context.Context, postID int64) (*models.PostDetail, error) {
// 	return l.postDao.GetDetail(ctx, postID)
// }

// // ListPosts returns a paginated list. When ordering by score, vote counts
// // are retrieved from Redis in the vote logic layer; here we use the DB fallback.
// func (l *PostLogic) ListPosts(ctx context.Context, params *models.PostListParams) ([]*models.Post, int64, error) {
// 	return l.postDao.List(ctx, params)
// }

// // DeletePost soft-deletes a post. Only the author can delete their own post.
// func (l *PostLogic) DeletePost(ctx context.Context, postID int64, requesterID int64) error {
// 	post, err := l.postDao.GetByID(ctx, postID)
// 	if err != nil {
// 		return err
// 	}
// 	if post.AuthorID != requesterID {
// 		return mysql.ErrUnauthorised
// 	}
// 	return l.postDao.Delete(ctx, postID)
// }

// // GetPostAge returns how long ago a post was created — useful for vote expiry rules.
// func (l *PostLogic) GetPostAge(ctx context.Context, postID int64) (time.Duration, error) {
// 	post, err := l.postDao.GetByID(ctx, postID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return time.Since(post.CreateTime), nil
// }

// // ErrPostTooOld is returned when a user tries to vote on an expired post.
// var ErrPostTooOld = errors.New("post is too old to vote on")

package logic

import (
	"context"
	"errors"
	"time"

	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// PostLogic handles business rules for posts.
type PostLogic struct {
	postDao mysql.PostStore
	voteDao mysql.PostVoteStore
}

// NewPostLogic wires up the logic layer.
// Accepts interfaces so tests can inject mocks.
func NewPostLogic(postDao mysql.PostStore, voteDao mysql.PostVoteStore) *PostLogic {
	return &PostLogic{
		postDao: postDao,
		voteDao: voteDao,
	}
}

// CreatePost validates input, generates a snowflake ID, and persists the post.
func (l *PostLogic) CreatePost(ctx context.Context, params *models.CreatePostParams) (*models.Post, error) {
	post := &models.Post{
		PostID:      snowflake.GenID(),
		Title:       params.Title,
		Content:     params.Content,
		AuthorID:    params.AuthorID,
		CommunityID: params.CommunityID,
		Status:      models.PostStatusNormal,
	}

	if err := l.postDao.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

// GetPostDetail fetches the full post with author, community and vote info.
func (l *PostLogic) GetPostDetail(ctx context.Context, postID int64) (*models.PostDetail, error) {
	return l.postDao.GetDetail(ctx, postID)
}

// ListPosts returns a paginated list.
func (l *PostLogic) ListPosts(ctx context.Context, params *models.PostListParams) ([]*models.Post, int64, error) {
	return l.postDao.List(ctx, params)
}

// DeletePost soft-deletes a post. Only the author can delete their own post.
func (l *PostLogic) DeletePost(ctx context.Context, postID int64, requesterID int64) error {
	post, err := l.postDao.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != requesterID {
		return mysql.ErrUnauthorised
	}
	return l.postDao.Delete(ctx, postID)
}

// GetPostAge returns how long ago a post was created — useful for vote expiry rules.
func (l *PostLogic) GetPostAge(ctx context.Context, postID int64) (time.Duration, error) {
	post, err := l.postDao.GetByID(ctx, postID)
	if err != nil {
		return 0, err
	}
	return time.Since(post.CreateTime), nil
}

// ErrPostTooOld is returned when a user tries to vote on an expired post.
var ErrPostTooOld = errors.New("post is too old to vote on")
