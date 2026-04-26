package controller

// import (
// 	"errors"
// 	"net/http"
// 	"repeat_practice/logic"
// 	"repeat_practice/models"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"
// )

// type UserController struct {
// 	userlogic *logic.UserLogic
// }

// func NewUserController(ul *logic.UserLogic) *UserController {
// 	return &UserController{userlogic: ul}
// }

// func (uc *UserController) SignUpHandler(ctx *gin.Context) {
// 	var p models.ParamSignUp
// 	if err := ctx.ShouldBindJSON(&p); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": formatValidationError(err)})
// 		return
// 	}

// 	if err := uc.userlogic.SignUp(&p); err != nil {
// 		// zap.L().Warn("user signup failed", zap.Error(err))预期错误就不打了
// 		// 日志放default

// 		switch {
// 		case errors.Is(err, logic.ErrUsernameTaken):
// 			ctx.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "Username is already taken"})

// 		case errors.Is(err, logic.ErrEmailTaken):
// 			ctx.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "Email is already registered"})

// 		default:
// 			// 生产环境：对用户友好 + 内部详细日志
// 			zap.L().Error("unexpected error in user signup", zap.Error(err))
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"code": 500,
// 				"msg":  "Registration failed, please try again later",
// 			})
// 		}
// 		return
// 	}

// 	// 注册成功，通常返回 201 Created，可带简单消息或用户ID
// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"code": 201,
// 		"msg":  "Registration successful.",
// 	})
// }

// // 登录
// func (uc *UserController) LoginHandler(ctx *gin.Context) {
// 	var p models.ParamLogin
// 	if err := ctx.ShouldBindJSON(&p); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"code": 400,
// 			"msg":  formatValidationError(err),
// 		})
// 		return
// 	}

// 	user, err := uc.userlogic.Login(&p)
// 	if err != nil {
// 		// 登录失败日志（注意这里可以打，因为是安全相关）
// 		zap.L().Warn("login failed",
// 			zap.String("email", p.Email),
// 			zap.Error(err),
// 		)

// 		switch {
// 		case errors.Is(err, logic.ErrAccountSuspended):
// 			ctx.JSON(http.StatusForbidden, gin.H{
// 				"code": 403,
// 				"msg":  err.Error(),
// 			})

// 		default:
// 			// 防止用户枚举攻击（用户名/邮箱探测）
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"code": 401,
// 				"msg":  "invalid email or password",
// 			})
// 		}
// 		return
// 	}

// 	// 登录成功
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"code": 200,
// 		"data": gin.H{
// 			"access_token": user.Token,
// 			"token_type":   "Bearer",
// 			"expires_in":   86400,
// 		},
// 	})
// }
