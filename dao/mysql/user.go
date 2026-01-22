package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

//把每一步数据库操作封装成函数
//待logic层根据业务需求调用

const secret = "ziying encryption"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	ePassword := encryptPassword(user.Password)
	//执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, ePassword)
	return
}

// md5基本不用了，
// bcrypt,scrypt,argon2现在公司用于加密比较多
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	inputPassword := user.Password /// 先保存用戶輸入的明文密碼
	var dbUser models.User
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(&dbUser, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库失败
		return err
	}
	// 比對密碼
	if encryptPassword(inputPassword) != dbUser.Password {
		return ErrorInvalidPassword
	}
	// 把 dbUser 的所有字段值，完整复制到 user 指向的结构体中
	*user = dbUser
	// 错误写法user = dbUser,
	// 不加*，只是让函数内部的局部变量 user（指针本身）改指向 dbUser，但调用方原来的结构体完全没变。
	// 函数返回后，外面的 u 还是只有 Username，其他字段还是零值或旧值。
	return nil
}
