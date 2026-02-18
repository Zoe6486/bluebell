package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/jmoiron/sqlx"
)

const secret = "ziying encryption"

var (
	ErrorUserNotExist    = errors.New("user not exist")
	ErrorInvalidPassword = errors.New("invalid password")
)

type UserDao struct {
	db *sqlx.DB
}

// 构造函数
func NewUserDao(db *sqlx.DB) *UserDao {
	return &UserDao{db: db}
}

func (u *UserDao) CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := u.db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return nil
}

func (u *UserDao) InsertUser(user *models.User) error {
	ePassword := encryptPassword(user.Password)
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err := u.db.Exec(sqlStr, user.UserID, user.Username, ePassword)
	return err
}

// md5基本不用了，
// bcrypt,scrypt,argon2现在公司用于加密比较多
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func (u *UserDao) Login(user *models.User) error {
	inputPassword := user.Password
	var dbUser models.User
	sqlStr := `select user_id, username, password from user where username=?`
	err := u.db.Get(&dbUser, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}

	if encryptPassword(inputPassword) != dbUser.Password {
		return ErrorInvalidPassword
	}

	*user = dbUser
	return nil
}
