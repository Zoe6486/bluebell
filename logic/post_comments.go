// package logic

// import (
// 	"context"
// 	"errors"

// 	"bluebell/dao/mysql"
// 	"bluebell/models"
// )

// // PostCommentLogic handles business rules for post comments.
// type PostCommentLogic struct {
// 	commentDao mysql.PostCommentStore
// 	postDao    mysql.PostStore
// }

// // NewPostCommentLogic wires up comment logic with its dependencies.
// func NewPostCommentLogic(commentDao mysql.PostCommentStore, postDao mysql.PostStore) *PostCommentLogic {
// 	return &PostCommentLogic{
// 		commentDao: commentDao,
// 		postDao:    postDao,
// 	}
// }

// // CreateComment validates that the parent post is active before inserting.
// func (l *PostCommentLogic) CreateComment(ctx context.Context, params *models.CreateCommentParams) (*models.PostComment, error) {
// 	// Verify the post exists and is not deleted/pending
// 	_, err := l.postDao.GetByID(ctx, int64(params.PostID))
// 	if err != nil {
// 		if errors.Is(err, mysql.ErrNotFound) {
// 			return nil, mysql.ErrNotFound
// 		}
// 		return nil, err
// 	}

// 	comment := &models.PostComment{
// 		PostID:  params.PostID,
// 		UserID:  params.UserID,
// 		Content: params.Content,
// 		Status:  models.CommentStatusNormal,
// 	}

// 	if err := l.commentDao.Create(ctx, comment); err != nil {
// 		return nil, err
// 	}
// 	return comment, nil
// }

// // ListComments returns a paginated list of comments for a post.
// func (l *PostCommentLogic) ListComments(ctx context.Context, params *models.CommentListParams) ([]*models.CommentDetail, int64, error) {
// 	return l.commentDao.ListByPost(ctx, params)
// }

// // DeleteComment allows the owner to soft-delete their comment.
//
//	func (l *PostCommentLogic) DeleteComment(ctx context.Context, commentID uint64, requesterID uint64) error {
//		return l.commentDao.Delete(ctx, commentID, requesterID)
//	}
package logic

import (
	"context"

	"bluebell/models"
)

type PostCommentStore interface {
	Create(ctx context.Context, c *models.PostComment) error
	GetByID(ctx context.Context, commentID uint64) (*models.PostComment, error)
	ListByPost(ctx context.Context, params *models.CommentListParams) ([]*models.CommentDetail, int64, error)
	Delete(ctx context.Context, commentID uint64, userID uint64) error
	IncrLikeCount(ctx context.Context, commentID uint64, delta int) error
}

type PostCommentLogic struct {
	commentStore PostCommentStore
	postStore    PostStore
}

func NewPostCommentLogic(commentStore PostCommentStore, postStore PostStore) *PostCommentLogic {
	return &PostCommentLogic{
		commentStore: commentStore,
		postStore:    postStore,
	}
}

func (l *PostCommentLogic) CreateComment(ctx context.Context, params *models.CreateCommentParams) (*models.PostComment, error) {
	_, err := l.postStore.GetByID(ctx, int64(params.PostID))
	if err != nil {
		return nil, err
	}

	comment := &models.PostComment{
		PostID:  params.PostID,
		UserID:  params.UserID,
		Content: params.Content,
		Status:  models.CommentStatusNormal,
	}

	if err := l.commentStore.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (l *PostCommentLogic) ListComments(ctx context.Context, params *models.CommentListParams) ([]*models.CommentDetail, int64, error) {
	return l.commentStore.ListByPost(ctx, params)
}

func (l *PostCommentLogic) DeleteComment(ctx context.Context, commentID uint64, requesterID uint64) error {
	return l.commentStore.Delete(ctx, commentID, requesterID)
}
