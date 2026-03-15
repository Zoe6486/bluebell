// package logic

// import (
// 	"bluebell/models"
// 	"bluebell/pkg/jwt"
// 	"bluebell/pkg/snowflake"
// )

// // 定义 interface，Logic 不依赖 mysql 包
// type UserStore interface {
// 	CheckUserExist(username string) error
// 	InsertUser(user *models.User) error
// 	Login(user *models.User) error
// }

// type UserLogic struct {
// 	store UserStore
// }

// // 构造函数
// func NewUserLogic(store UserStore) *UserLogic {
// 	return &UserLogic{store: store}
// }
// func (l *UserLogic) SignUp(p *models.ParamSignUp) error {
// 	// 判断用户是否存在
// 	if err := l.store.CheckUserExist(p.Username); err != nil {
// 		return err
// 	}

// 	// 生成 uid
// 	userID := snowflake.GenID()

// 	user := &models.User{
// 		UserID:   userID,
// 		Username: p.Username,
// 		Password: p.Password,
// 	}

// 	// 保存进数据库
// 	if err := l.store.InsertUser(user); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l *UserLogic) Login(p *models.ParamLogin) (*models.User, error) {
// 	user := &models.User{
// 		Username: p.Username,
// 		Password: p.Password,
// 	}

// 	// 调用 dao 登录
// 	if err := l.store.Login(user); err != nil {
// 		return nil, err
// 	}

// 	// 生成 JWT
// 	token, err := jwt.GenToken(user.UserID, user.Username)
// 	if err != nil {
// 		return nil, err
// 	}

//		user.Token = token
//		return user, nil
//	}
package logic

import (
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserLogic struct {
	store UserStore
}

type UserStore interface {
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Insert(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByUserID(userID int64) (*models.User, error)
}

func NewUserLogic(store UserStore) *UserLogic {
	return &UserLogic{store: store}
}

func (l *UserLogic) SignUp(p *models.ParamSignUp) error {
	if taken, err := l.store.ExistsByUsername(p.Username); err != nil {
		return err
	} else if taken {
		return ErrUsernameTaken
	}

	if taken, err := l.store.ExistsByEmail(p.Email); err != nil {
		return err
	} else if taken {
		return ErrEmailTaken
	}

	user := &models.User{
		UserID:   snowflake.GenID(),
		Username: p.Username,
		Email:    p.Email,
		Password: p.Password,
	}

	return l.store.Insert(user)
}

func (l *UserLogic) Login(p *models.ParamLogin) (*models.User, error) {
	user, err := l.store.GetByEmail(p.Email)
	if err != nil {
		return nil, ErrInvalidCredentials // don't leak whether email exists
	}

	if user.Status == models.UserStatusSuspended {
		return nil, ErrAccountSuspended
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return user, nil
}

// Sentinel errors for the logic layer
var (
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountSuspended   = errors.New("account suspended")
)
