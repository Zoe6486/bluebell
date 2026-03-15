// package mysql

// import (
// 	"bluebell/models"
// 	"crypto/md5"
// 	"database/sql"
// 	"encoding/hex"
// 	"errors"

// 	"github.com/jmoiron/sqlx"
// )

// const secret = "ziying encryption"

// var (
// 	ErrorUserNotExist    = errors.New("user not exist")
// 	ErrorInvalidPassword = errors.New("invalid password")
// )

// type UserDao struct {
// 	db *sqlx.DB
// }

// // 构造函数
// func NewUserDao(db *sqlx.DB) *UserDao {
// 	return &UserDao{db: db}
// }

// func (u *UserDao) CheckUserExist(username string) error {
// 	sqlStr := `select count(user_id) from user where username = ?`
// 	var count int64
// 	if err := u.db.Get(&count, sqlStr, username); err != nil {
// 		return err
// 	}
// 	if count > 0 {
// 		return errors.New("用户已存在")
// 	}
// 	return nil
// }

// func (u *UserDao) InsertUser(user *models.User) error {
// 	ePassword := encryptPassword(user.Password)
// 	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
// 	_, err := u.db.Exec(sqlStr, user.UserID, user.Username, ePassword)
// 	return err
// }

// // md5基本不用了，
// // bcrypt,scrypt,argon2现在公司用于加密比较多
// func encryptPassword(oPassword string) string {
// 	h := md5.New()
// 	h.Write([]byte(secret))
// 	return hex.EncodeToString(h.Sum([]byte(oPassword)))
// }

// func (u *UserDao) Login(user *models.User) error {
// 	inputPassword := user.Password
// 	var dbUser models.User
// 	sqlStr := `select user_id, username, password from user where username=?`
// 	err := u.db.Get(&dbUser, sqlStr, user.Username)
// 	if err == sql.ErrNoRows {
// 		return ErrorUserNotExist
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	if encryptPassword(inputPassword) != dbUser.Password {
// 		return ErrorInvalidPassword
// 	}

//		*user = dbUser
//		return nil
//	}
package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUsernameTaken    = errors.New("username already taken")
	ErrEmailTaken       = errors.New("email already registered")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrAccountSuspended = errors.New("account suspended")
)

type UserStore interface {
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Insert(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByUserID(userID int64) (*models.User, error)
}

type userDao struct {
	db *sqlx.DB
}

func NewUserDao(db *sqlx.DB) UserStore {
	return &userDao{db: db}
}

func (d *userDao) ExistsByUsername(username string) (bool, error) {
	var count int
	err := d.db.Get(&count, `SELECT COUNT(1) FROM user WHERE username = ?`, username)
	return count > 0, err
}

func (d *userDao) ExistsByEmail(email string) (bool, error) {
	var count int
	err := d.db.Get(&count, `SELECT COUNT(1) FROM user WHERE email = ?`, email)
	return count > 0, err
}

func (d *userDao) Insert(user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(
		`INSERT INTO user (user_id, username, email, password) VALUES (?, ?, ?, ?)`,
		user.UserID, user.Username, user.Email, string(hash),
	)
	return err
}

func (d *userDao) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := d.db.Get(&u,
		`SELECT id, user_id, username, email, password, status, created_at, updated_at
		 FROM user WHERE email = ? AND status != ?`,
		email, models.UserStatusDeleted,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &u, err
}

func (d *userDao) GetByUserID(userID int64) (*models.User, error) {
	var u models.User
	err := d.db.Get(&u,
		`SELECT id, user_id, username, email, password, status, created_at, updated_at
		 FROM user WHERE user_id = ? AND status != ?`,
		userID, models.UserStatusDeleted,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &u, err
}
