// package controller

// import (
// 	"errors"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"

// 	"bluebell/logic"
// 	"bluebell/models"
// )

// type PostController struct {
// 	postLogic *logic.PostLogic
// }

// func NewPostController(postLogic *logic.PostLogic) *PostController {
// 	return &PostController{postLogic: postLogic}
// }

// func (c *PostController) CreatePost(ctx *gin.Context) {
// 	var params models.CreatePostParams
// 	if err := ctx.ShouldBindJSON(&params); err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, formatValidationError(err))
// 		return
// 	}

// 	params.AuthorID = getAuthorizedUserID(ctx)

// 	post, err := c.postLogic.CreatePost(ctx.Request.Context(), &params)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
// 		return
// 	}

// 	ResponseCreated(ctx, post)
// }

// func (c *PostController) GetPost(ctx *gin.Context) {
// 	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
// 		return
// 	}

// 	detail, err := c.postLogic.GetPostDetail(ctx.Request.Context(), postID)
// 	if err != nil {
// 		if errors.Is(err, logic.ErrNotFound) {
// 			ResponseError(ctx, http.StatusNotFound, "post not found")
// 			return
// 		}
// 		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
// 		return
// 	}

// 	ResponseSuccess(ctx, detail)
// }

// func (c *PostController) ListPosts(ctx *gin.Context) {
// 	var params models.PostListParams
// 	if err := ctx.ShouldBindQuery(&params); err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, formatValidationError(err))
// 		return
// 	}

// 	posts, total, err := c.postLogic.ListPosts(ctx.Request.Context(), &params)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
// 		return
// 	}

// 	ResponseSuccess(ctx, gin.H{
// 		"list":      posts,
// 		"total":     total,
// 		"page":      params.Page,
// 		"page_size": params.PageSize,
// 	})
// }

// func (c *PostController) DeletePost(ctx *gin.Context) {
// 	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
// 	if err != nil {
// 		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
// 		return
// 	}

// 	if err := c.postLogic.DeletePost(ctx.Request.Context(), postID, getAuthorizedUserID(ctx)); err != nil {
// 		switch {
// 		case errors.Is(err, logic.ErrNotFound):
// 			ResponseError(ctx, http.StatusNotFound, "post not found")
// 		case errors.Is(err, logic.ErrUnauthorised):
// 			ResponseError(ctx, http.StatusForbidden, "you can only delete your own posts")
// 		default:
// 			ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
// 		}
// 		return
// 	}

// 	ResponseNoContent(ctx)
// }

package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bluebell/logic"
	"bluebell/models"
)

type PostController struct {
	postLogic *logic.PostLogic
}

func NewPostController(postLogic *logic.PostLogic) *PostController {
	return &PostController{postLogic: postLogic}
}

func (c *PostController) CreatePost(ctx *gin.Context) {
	var params models.CreatePostParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, formatValidationError(err))
		return
	}

	params.AuthorID = getAuthorizedUserID(ctx)

	post, err := c.postLogic.CreatePost(ctx.Request.Context(), &params)
	if err != nil {
		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		return
	}

	ResponseCreated(ctx, post)
}

func (c *PostController) GetPost(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
		return
	}

	detail, err := c.postLogic.GetPostDetail(ctx.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			ResponseError(ctx, http.StatusNotFound, "post not found")
			return
		}
		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		return
	}

	ResponseSuccess(ctx, detail)
}

func (c *PostController) ListPosts(ctx *gin.Context) {
	var params models.PostListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, formatValidationError(err))
		return
	}

	posts, total, err := c.postLogic.ListPosts(ctx.Request.Context(), &params)
	if err != nil {
		ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		return
	}

	ResponseSuccess(ctx, gin.H{
		"list":      posts,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
	})
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
		return
	}

	if err := c.postLogic.DeletePost(ctx.Request.Context(), postID, getAuthorizedUserID(ctx)); err != nil {
		switch {
		case errors.Is(err, logic.ErrNotFound):
			ResponseError(ctx, http.StatusNotFound, "post not found")
		case errors.Is(err, logic.ErrUnauthorised):
			ResponseError(ctx, http.StatusForbidden, "you can only delete your own posts")
		default:
			ResponseError(ctx, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	ResponseNoContent(ctx)
}
