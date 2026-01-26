package controller

import (
	"bluebell/middleware"

	"github.com/gin-gonic/gin"
)

// GetUserID 从 context 里安全获取 userID，失败时返回 0 和 false
func GetUserID(c *gin.Context) (int64, bool) {
	return middleware.CurrentUserID(c)
}

// // BindAndValidate 绑定并校验请求参数，统一错误返回
// func BindAndValidate(c *gin.Context, obj interface{}) bool {
// 	if err := c.ShouldBindJSON(obj); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return false
// 	}
// 	// 这里可以加你额外的校验逻辑，比如 validator 校验
// 	return true
// }
