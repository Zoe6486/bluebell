package controller

import (
	//"bluebell/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化自定义参数校验逻辑
func InitValidator() (err error) {
	// binding.Validator.Engine() 返回 Gin 默认使用的 validator 引擎
	// 类型断言(Type Assertion)成 *validator.Validate，才能注册自定义逻辑
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		// 直接给命名返回值 err 赋值，然后 return
		err = fmt.Errorf("gin validator engine is not *validator.Validate")
		return
	}
	// -----------------------------
	// 1. 自定义字段名显示规则
	// -----------------------------
	// 默认 validator 报错信息使用的是 struct 的字段名（如 Password）
	// 这里改成使用 JSON tag（如 "password"），便于前端理解
	// 	{
	//   "Username": "Username is a required field",
	//   "Password": "Password is a required field"
	// }改成
	// 	{
	//   "username": "username is a required field",
	//   "password": "password is a required field"
	// }

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	//可以不写，validator有自带的
	// 	如果业务逻辑复杂，单靠 tag 校验不够，就会写 RegisterStructValidation

	// 例如：

	// “开始时间必须早于结束时间”

	// “用户等级为 VIP 才能填写某个字段”

	// 这类接口才会用你 validator.go 里那种结构体级校验函数
	// -----------------------------
	// 2. 注册结构体级别校验
	// -----------------------------
	// 结构体级校验可以同时检查多个字段之间的逻辑关系
	// 这里给 ParamSignUp 注册了一个自定义校验函数
	//v.RegisterStructValidation(SignUpParamStructLevelValidation, models.ParamSignUp{})
	// 如果一切正常，err 默认为 nil，会被返回
	return
}

// func SignUpParamStructLevelValidation(sl validator.StructLevel) {
// 	// sl.Current() 返回当前正在校验的结构体对象
// 	// Interface() 获取其接口类型，再断言为 ParamSignUp
// 	// 如果两次密码不一致，就向 validator 报错
// 	// 参数说明：
// 	// 1. su.RePassword -> 实际出错的字段值
// 	// 2. "re_password" -> json 名称
// 	// 3. "RePassword" -> 原字段名（结构体字段）
// 	// 4. "eqfield" -> 校验规则标识，这里表示要等于另一个字段
// 	// 5. "password" -> 要比较的字段名称
// 	su := sl.Current().Interface().(models.ParamSignUp)
// 	if su.Password != su.RePassword {
// 		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")
// 	}
// }

// removeTopStruct 去除提示信息中的结构体名称
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
