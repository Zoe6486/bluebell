package logic

import (
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 定义 interface，Logic 不依赖 mysql 包
type UserStore interface {
	CheckUserExist(username string) error
	InsertUser(user *models.User) error
	Login(user *models.User) error
}

type UserLogic struct {
	store UserStore
}

// 构造函数
func NewUserLogic(store UserStore) *UserLogic {
	return &UserLogic{store: store}
}
func (l *UserLogic) SignUp(p *models.ParamSignUp) error {
	// 判断用户是否存在
	if err := l.store.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 生成 uid
	userID := snowflake.GenID()

	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 保存进数据库
	if err := l.store.InsertUser(user); err != nil {
		return err
	}

	return nil
}

func (l *UserLogic) Login(p *models.ParamLogin) (*models.User, error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 调用 dao 登录
	if err := l.store.Login(user); err != nil {
		return nil, err
	}

	// 生成 JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return user, nil
}
