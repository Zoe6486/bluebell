package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	logic *logic.CommunityLogic
}

func NewCommunityController(logic *logic.CommunityLogic) *CommunityController {
	return &CommunityController{logic: logic}
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// 参数错误，给 400 Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid community id"})
		return
	}

	data, err := c.logic.GetCommunityDetail(id)
	if err != nil {
		//  重点：根据 DAO 返回的具体错误给状态码
		if err == mysql.ErrorInvalidID {
			// 没找到资源，给 404 Not Found
			ctx.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//  直接返回对象
	ctx.JSON(http.StatusOK, data)
}
