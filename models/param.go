package models

// // 定义请求的参数结构体
// // ParamSignUp 注册请求参数
// type ParamSignUp struct {
// 	Username   string `json:"username" binding:"required"`
// 	Password   string `json:"password" binding:"required"`
// 	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
// }

// // ParamLogin 登录请求参数
//
//	type ParamLogin struct {
//		Username string `json:"username" binding:"required"`
//		Password string `json:"password" binding:"required"`
//	}
type ParamSignUp struct {
	Username        string `json:"username"         binding:"required,min=3,max=32"`
	Email           string `json:"email"            binding:"required,email"`
	Password        string `json:"password"         binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
