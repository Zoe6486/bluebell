package controller

import (
	"bluebell/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard API response envelope
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ResponseSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created successfully, 201
func ResponseCreated(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Deleted successfully, 204 No Content
func ResponseNoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func ResponseError(ctx *gin.Context, httpStatus int, msg string) {
	ctx.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: msg,
	})
}

// getAuthorizedUserID retrieves the user ID injected by the JWT middleware.
// The middleware should store it under the key "user_id".
func getAuthorizedUserID(ctx *gin.Context) int64 {
	// userID, _ := ctx.Get("user_id")
	// id, _ := userID.(int64)
	// return id
	id, _ := middleware.CurrentUserID(ctx)
	return id
}
