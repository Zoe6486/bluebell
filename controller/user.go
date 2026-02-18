package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserController 包含 UserLogic 的依赖
type UserController struct {
	logic *logic.UserLogic
}

// 构造函数
func NewUserController(l *logic.UserLogic) *UserController {
	return &UserController{logic: l}
}

func (uc *UserController) SignUpHandler(c *gin.Context) {
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			errsMap := make(map[string]string)
			for _, e := range errs {
				errsMap[e.Namespace()] = "字段验证失败: " + e.Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errsMap)})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用注入的 logic
	if err := uc.logic.SignUp(&p); err != nil {
		zap.L().Error("UserLogic.SignUp failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "sign up successfully"})
}

func (uc *UserController) LoginHandler(c *gin.Context) {
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		if errs, ok := err.(validator.ValidationErrors); ok {
			errsMap := make(map[string]string)
			for _, e := range errs {
				errsMap[e.Namespace()] = "字段验证失败: " + e.Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errsMap)})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.logic.Login(&p)
	if err != nil {
		zap.L().Error("UserLogic.Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": user.Token,
		"token_type":   "Bearer",
		"expires_in":   86400,
	})
}
