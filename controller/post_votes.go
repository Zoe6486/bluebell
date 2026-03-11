package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"bluebell/logic"
	"bluebell/models"
)

// PostVoteController handles HTTP for post votes.
type PostVoteController struct {
	voteLogic *logic.PostVoteLogic
}

// NewPostVoteController wires up the controller.
func NewPostVoteController(voteLogic *logic.PostVoteLogic) *PostVoteController {
	return &PostVoteController{voteLogic: voteLogic}
}

// RegisterRoutes attaches vote routes to a router group.
func (c *PostVoteController) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/vote", c.VotePost)
}

// VotePost godoc
// POST /api/v1/posts/vote
// Body: { "post_id": 123, "vote_type": 1 }
func (c *PostVoteController) VotePost(ctx *gin.Context) {
	var params models.VotePostParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	params.UserID = getAuthorizedUserID(ctx)

	if err := c.voteLogic.VotePost(ctx.Request.Context(), &params); err != nil {
		switch {
		case errors.Is(err, logic.ErrPostTooOld):
			ResponseError(ctx, http.StatusUnprocessableEntity, "this post is too old to vote on")
		default:
			ResponseError(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ResponseSuccess(ctx, nil)
}
