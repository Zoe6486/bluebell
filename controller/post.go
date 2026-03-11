package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
)

// PostController handles HTTP for posts.
type PostController struct {
	postLogic *logic.PostLogic
}

// NewPostController wires up the controller.
func NewPostController(postLogic *logic.PostLogic) *PostController {
	return &PostController{postLogic: postLogic}
}

// RegisterRoutes attaches all post routes to a router group.
func (c *PostController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("", c.CreatePost)
	r.GET("", c.ListPosts)
	r.GET("/:post_id", c.GetPost)
	r.DELETE("/:post_id", c.DeletePost)
}

// CreatePost godoc
// POST /api/v1/posts
func (c *PostController) CreatePost(ctx *gin.Context) {
	var params models.CreatePostParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Inject authenticated user ID from JWT middleware
	params.AuthorID = getAuthorizedUserID(ctx)

	post, err := c.postLogic.CreatePost(ctx.Request.Context(), &params)
	if err != nil {
		ResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseCreated(ctx, post)
}

// GetPost godoc
// GET /api/v1/posts/:post_id
func (c *PostController) GetPost(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
		return
	}

	detail, err := c.postLogic.GetPostDetail(ctx.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, mysql.ErrNotFound) {
			ResponseError(ctx, http.StatusNotFound, "post not found")
			return
		}
		ResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(ctx, detail)
}

// ListPosts godoc
// GET /api/v1/posts?page=1&page_size=20&community_id=1&order_by=time
func (c *PostController) ListPosts(ctx *gin.Context) {
	var params models.PostListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	posts, total, err := c.postLogic.ListPosts(ctx.Request.Context(), &params)
	if err != nil {
		ResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(ctx, gin.H{
		"list":  posts,
		"total": total,
		"page":  params.Page,
	})
}

// DeletePost godoc
// DELETE /api/v1/posts/:post_id
func (c *PostController) DeletePost(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		ResponseError(ctx, http.StatusBadRequest, "invalid post_id")
		return
	}

	userID := getAuthorizedUserID(ctx)

	if err := c.postLogic.DeletePost(ctx.Request.Context(), postID, userID); err != nil {
		switch {
		case errors.Is(err, mysql.ErrNotFound):
			ResponseError(ctx, http.StatusNotFound, "post not found")
		case errors.Is(err, mysql.ErrUnauthorised):
			ResponseError(ctx, http.StatusForbidden, "you can only delete your own posts")
		default:
			ResponseError(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ResponseNoContent(ctx)
}
