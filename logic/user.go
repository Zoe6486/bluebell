package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
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

func Login(p *models.ParamLogin) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
