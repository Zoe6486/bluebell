// package controller

// import (
// 	"bluebell/logic"
// 	"bluebell/models"
// 	"bluebell/pkg/response"
// 	"errors"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"
// )

// type UserController struct {
// 	userlogic *logic.UserLogic
// }

// func NewUserController(ul *logic.UserLogic) *UserController {
// 	return &UserController{userlogic: ul}
// }

// // SignUpHandler 注册
// func (uc *UserController) SignUpHandler(ctx *gin.Context) {
// 	var p models.ParamSignUp

// 	if err := ctx.ShouldBindJSON(&p); err != nil {
// 		response.ValidationError(ctx, formatValidationError(err))
// 		return
// 	}

// 	if err := uc.userlogic.SignUp(&p); err != nil {
// 		switch {
// 		case errors.Is(err, logic.ErrUsernameTaken):
// 			response.Error(ctx, http.StatusConflict, "username already taken")

// 		case errors.Is(err, logic.ErrEmailTaken):
// 			response.Error(ctx, http.StatusConflict, "email already registered")

// 		default:
// 			zap.L().Error("unexpected error in user signup", zap.Error(err))
// 			response.Error(ctx, http.StatusInternalServerError, "registration failed")
// 		}
// 		return
// 	}

// 	// 201 Created（RESTful：成功不包 code/msg）
// 	response.Success(ctx, http.StatusCreated, gin.H{
// 		"message": "registration successful",
// 	})
// }

// // LoginHandler 登录
// func (uc *UserController) LoginHandler(ctx *gin.Context) {
// 	var p models.ParamLogin

// 	if err := ctx.ShouldBindJSON(&p); err != nil {
// 		response.ValidationError(ctx, formatValidationError(err))
// 		return
// 	}

// 	user, err := uc.userlogic.Login(&p)
// 	if err != nil {
// 		zap.L().Warn("login failed",
// 			zap.String("email", p.Email),
// 			zap.Error(err),
// 		)

// 		switch {
// 		case errors.Is(err, logic.ErrAccountSuspended):
// 			response.Error(ctx, http.StatusForbidden, err.Error())

// 		default:
// 			// RESTful + 安全最佳实践（不暴露具体原因）
// 			response.Error(ctx, http.StatusUnauthorized, "invalid email or password")
// 		}
// 		return
// 	}

// 	// 登录成功（直接返回资源）
// 	response.Success(ctx, http.StatusOK, gin.H{
// 		"access_token": user.Token,
// 		"token_type":   "Bearer",
// 		"expires_in":   86400,
// 	})
// }

package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	userlogic *logic.UserLogic
}

func NewUserController(ul *logic.UserLogic) *UserController {
	return &UserController{userlogic: ul}
}

// SignUpHandler 注册
func (uc *UserController) SignUpHandler(ctx *gin.Context) {
	var p models.ParamSignUp

	// 1️⃣ 参数绑定 + 校验错误
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errs := FormatValidationError(err)
		response.ValidationError(ctx, errs)
		return
	}

	// 2️⃣ 业务错误处理
	if err := uc.userlogic.SignUp(&p); err != nil {
		switch {
		case errors.Is(err, logic.ErrUsernameTaken):
			response.Conflict(ctx, "username already registered")

		case errors.Is(err, logic.ErrEmailTaken):
			response.Conflict(ctx, "email already registered")

		case errors.Is(err, logic.ErrDuplicateEntry):
			// 👉 兜底（不暴露具体字段）
			response.Conflict(ctx, "resource already exists")

		default:
			zap.L().Error("unexpected error in signup", zap.Error(err))
			response.InternalError(ctx, err)
		}
		return
	}

	// 3️⃣ 成功
	response.Created(ctx, gin.H{
		"message": "registration successful",
	})
}

// LoginHandler 登录
func (uc *UserController) LoginHandler(ctx *gin.Context) {
	var p models.ParamLogin

	// 1️⃣ 参数校验
	if err := ctx.ShouldBindJSON(&p); err != nil {
		errs := FormatValidationError(err)
		response.ValidationError(ctx, errs)
		return
	}

	// 2️⃣ 业务逻辑
	user, err := uc.userlogic.Login(&p)
	if err != nil {
		// 🔒 安全日志（不暴露给用户）
		switch {
		case errors.Is(err, logic.ErrAccountSuspended):
			zap.L().Warn("login attempt on suspended account",
				zap.String("email", p.Email),
				zap.Error(err),
			)
		default:
			zap.L().Warn("invalid login attempt",
				zap.String("email", p.Email),
				zap.Error(err),
			)
		}

		// ❗统一返回（防止账号枚举）
		response.Unauthorized(ctx, "invalid email or password")
		return
	}

	// 3️⃣ 成功
	response.OK(ctx, gin.H{
		"access_token": user.Token,
		"token_type":   "Bearer",
		"expires_in":   86400, // "expires_in": jwt.TokenExpireSeconds后期写成这样
	})
}
