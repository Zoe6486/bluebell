package controller

import (
	"bluebell/logic"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	// Controller 依赖 Logic
	logic *logic.CommunityLogic
}

func NewCommunityController(l *logic.CommunityLogic) *CommunityController {
	return &CommunityController{logic: l}
}

// GetCommunityListHandler 列表接口
func (c *CommunityController) GetCommunityListHandler(ctx *gin.Context) {
	data, err := c.logic.GetCommunityList()
	if err != nil {
		//ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "获取列表失败"})
		// 服务器内部错误，直接给 500
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetCommunityDetailHandler 获取详情
func (c *CommunityController) GetCommunityDetailHandler(ctx *gin.Context) {
	// 获取 URL 中的参数 /community/:id
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// 参数错误，给 400 Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community id"})
		return
	}
	// 调用 Logic 层
	data, err := c.logic.GetCommunityDetail(id)
	if err != nil {
		// 只判断业务错误
		if errors.Is(err, logic.ErrCommunityNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		// 其他错误直接 500
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//  成功返回 (直接返回 models.CommunityDetail 对象)
	ctx.JSON(http.StatusOK, data)
}
