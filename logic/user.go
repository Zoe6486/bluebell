package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"fmt"
)

// 存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	fmt.Println("CheckUserExist:", p.Username)
	if err := mysql.CheckUserExist(p.Username); err != nil {
		fmt.Println("CheckUserExist failed:", err)
		return err
	}
	//生成uid
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//保存进数据库
	if err := mysql.InsertUser(user); err != nil {
		fmt.Println("InsertUser failed:", err)
		return err
	}
	fmt.Println("User inserted:", user.Username)
	return nil
	//redis.xxx ...
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT的token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
