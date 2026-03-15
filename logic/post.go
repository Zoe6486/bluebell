package logic

import (
	"context"
	"time"

	"bluebell/models"
	"bluebell/pkg/snowflake"
)

type PostStore interface {
	Create(ctx context.Context, p *models.Post) error
	GetByID(ctx context.Context, postID int64) (*models.Post, error)
	GetDetail(ctx context.Context, postID int64) (*models.PostDetail, error)
	List(ctx context.Context, params *models.PostListParams) ([]*models.Post, int64, error)
	Delete(ctx context.Context, postID int64) error
}

type PostLogic struct {
	postStore PostStore
}

func NewPostLogic(postStore PostStore) *PostLogic {
	return &PostLogic{postStore: postStore}
}

func (l *PostLogic) CreatePost(ctx context.Context, params *models.CreatePostParams) (*models.Post, error) {
	post := &models.Post{
		PostID:      snowflake.GenID(),
		Title:       params.Title,
		Content:     params.Content,
		AuthorID:    params.AuthorID,
		CommunityID: params.CommunityID,
		Status:      models.PostStatusNormal,
	}
	if err := l.postStore.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (l *PostLogic) GetPostDetail(ctx context.Context, postID int64) (*models.PostDetail, error) {
	return l.postStore.GetDetail(ctx, postID)
}

func (l *PostLogic) ListPosts(ctx context.Context, params *models.PostListParams) ([]*models.Post, int64, error) {
	return l.postStore.List(ctx, params)
}

func (l *PostLogic) DeletePost(ctx context.Context, postID int64, requesterID int64) error {
	post, err := l.postStore.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != requesterID {
		return ErrUnauthorised
	}
	return l.postStore.Delete(ctx, postID)
}

func (l *PostLogic) GetPostAge(ctx context.Context, postID int64) (time.Duration, error) {
	post, err := l.postStore.GetByID(ctx, postID)
	if err != nil {
		return 0, err
	}
	return time.Since(post.CreatedAt), nil
}
