// package controller

// import (
// 	"errors"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"

// 	"bluebell/dao/mysql"
// 	"bluebell/logic"
// 	"bluebell/models"
// )

// // PostCommentController handles HTTP for post comments.
// type PostCommentController struct {
// 	commentLogic *logic.PostCommentLogic
// }

// // NewPostCommentController wires up the controller.
// func NewPostCommentController(commentLogic *logic.PostCommentLogic) *PostCommentController {
// 	return &PostCommentController{commentLogic: commentLogic}
// }

// // RegisterRoutes attaches comment routes.
// // Typically nested under /posts/:post_id/comments
// func (c *PostCommentController) RegisterRoutes(r *gin.RouterGroup) {
// 	r.POST("", c.CreateComment)
// 	r.GET("", c.ListComments)
// 	r.DELETE("/:comment_id", c.DeleteComment)
// }

// // CreateComment godoc
// // POST /api/v1/posts/:post_id/comments
// func (c *PostCommentController) CreateComment(ctx *gin.Context) {
// 	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 64)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
// 		return
// 	}

// 	var params models.CreateCommentParams
// 	if err := ctx.ShouldBindJSON(&params); err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	params.PostID = postID
// 	params.UserID = uint64(getAuthorizedUserID(ctx))

// 	comment, err := c.commentLogic.CreateComment(ctx.Request.Context(), &params)
// 	if err != nil {
// 		if errors.Is(err, mysql.ErrNotFound) {
// 			ResponseError(ctx, http.StatusNotFound, "post not found")
// 			return
// 		}
// 		ResponseError(ctx, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	ResponseSuccess(ctx, comment)
// }

// // ListComments godoc
// // GET /api/v1/posts/:post_id/comments?page=1&page_size=20
// func (c *PostCommentController) ListComments(ctx *gin.Context) {
// 	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 64)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
// 		return
// 	}

// 	var params models.CommentListParams
// 	if err := ctx.ShouldBindQuery(&params); err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	params.PostID = postID

// 	comments, total, err := c.commentLogic.ListComments(ctx.Request.Context(), &params)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	ResponseSuccess(ctx, gin.H{
// 		"list":  comments,
// 		"total": total,
// 		"page":  params.Page,
// 	})
// }

// // DeleteComment godoc
// // DELETE /api/v1/posts/:post_id/comments/:comment_id
// func (c *PostCommentController) DeleteComment(ctx *gin.Context) {
// 	commentID, err := strconv.ParseUint(ctx.Param("comment_id"), 10, 64)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, "invalid comment_id")
// 		return
// 	}

// 	userID := uint64(getAuthorizedUserID(ctx))

// 	if err := c.commentLogic.DeleteComment(ctx.Request.Context(), commentID, userID); err != nil {
// 		if errors.Is(err, mysql.ErrNotFound) {
// 			ResponseError(ctx, http.StatusNotFound, "comment not found or already deleted")
// 			return
// 		}
// 		ResponseError(ctx, http.StatusInternalServerError, err.Error())
// 		return
// 	}

//		ResponseSuccess(ctx, nil)
//	}
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bluebell/logic"
	"bluebell/models"
)

type PostCommentController struct {
	commentLogic *logic.PostCommentLogic
}

func NewPostCommentController(commentLogic *logic.PostCommentLogic) *PostCommentController {
	return &PostCommentController{commentLogic: commentLogic}
}

func (c *PostCommentController) CreateComment(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
		return
	}

	var params models.CreateCommentParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, formatValidationError(err))
		return
	}

	params.PostID = postID
	params.UserID = uint64(getAuthorizedUserID(ctx))

	comment, err := c.commentLogic.CreateComment(ctx.Request.Context(), &params)
	if err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			ResponseError(ctx, http.StatusNotFound, "post not found")
			return
		}
		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		return
	}

	ResponseCreated(ctx, comment)
}

func (c *PostCommentController) ListComments(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
		return
	}

	var params models.CommentListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, formatValidationError(err))
		return
	}
	params.PostID = postID

	comments, total, err := c.commentLogic.ListComments(ctx.Request.Context(), &params)
	if err != nil {
		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		return
	}

	ResponseSuccess(ctx, gin.H{
		"list":      comments,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
	})
}

func (c *PostCommentController) DeleteComment(ctx *gin.Context) {
	commentID, err := strconv.ParseUint(ctx.Param("comment_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid comment_id")
		return
	}

	if err := c.commentLogic.DeleteComment(ctx.Request.Context(), commentID, uint64(getAuthorizedUserID(ctx))); err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			ResponseError(ctx, http.StatusNotFound, "comment not found or already deleted")
			return
		}
		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		return
	}

	ResponseNoContent(ctx)
}
