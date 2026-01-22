package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	//var p models.ParamSignUp

	p := new(models.ParamSignUp)
	//if err := c.ShouldBindJSON(&p); err != nil {
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断是否为 validator 校验错误
		if errs, ok := err.(validator.ValidationErrors); ok {
			// 格式化错误，返回前端友好 JSON
			errsMap := make(map[string]string)
			for _, e := range errs {
				// e.Field() 返回 struct 字段名
				// e.Error() 返回默认错误信息
				////错误信息如下：
				////	"error": {
				////	    "re_password": "Key: 'ParamSignUp.re_password' Error:Field validation for 're_password' failed on the 'required' tag"
				////	}
				//errsMap[e.Field()] = e.Error()
				//
				// 直接用 e.Namespace() 或 e.Field()，最终交给 removeTopStruct 处理

				errsMap[e.Namespace()] = "字段验证失败: " + e.Tag()
			}
			// 错误信息如下
			// 				{
			//     "error": {
			//         "ParamSignUp.re_password": "字段验证失败: required"
			//     }
			// }
			//c.JSON(http.StatusBadRequest, gin.H{"error": errsMap})
			// 如果想去掉结构体前缀，可以用 removeTopStruct
			// 				{
			//     "error": {
			//         "re_password": "字段验证失败: required"
			//     }
			//去掉了前面的 `ParamSignUp.`
			c.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errsMap)})
			return
		}

		// 其他错误直接返回
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// //手动对请求参数进行的业务规则校验（写了上面的下面就不用手动了）
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	// 	zap.L().Error("SignUp with invalid param")
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }
	// fmt.Println(p)

	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//3.返回响应
	c.JSON(http.StatusCreated, gin.H{
		"msg": "sign up successfully",
	})

}

// LogInHandler 登录
func LoginHandler(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断是否为 validator 校验错误
		if errs, ok := err.(validator.ValidationErrors); ok {
			// 格式化错误，返回前端友好 JSON
			errsMap := make(map[string]string)
			for _, e := range errs {
				// 直接用 e.Namespace() 或 e.Field()，最终交给 removeTopStruct 处理
				errsMap[e.Namespace()] = "字段验证失败: " + e.Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errsMap)})
			return
		}

		// 其他错误直接返回
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2.业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "username or password is not correct.",
		})
		return
	}
	// 3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "log in successfully.",
	})

}
