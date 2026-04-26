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
	"bluebell/logic"
	"bluebell/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// type UserStore interface {
// 	ExistsByUsername(username string) (bool, error)
// 	ExistsByEmail(email string) (bool, error)
// 	Insert(user *models.User) error
// 	GetByEmail(email string) (*models.User, error)
// 	GetByUserID(userID int64) (*models.User, error)
// }

type userDao struct {
	db *sqlx.DB
}

// ✅ 编译期检查
// “强制要求 userDao 必须实现 UserStore 接口
var _ logic.UserStore = (*userDao)(nil)

func NewUserDao(db *sqlx.DB) logic.UserStore {
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
	// hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }
	// _, err = d.db.Exec(
	// 	`INSERT INTO user (user_id, username, email, password) VALUES (?, ?, ?, ?)`,
	// 	user.UserID, user.Username, user.Email, string(hash),
	// )
	// return err
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. insert
	_, err = d.db.Exec(
		`INSERT INTO user(user_id, username, email, password) VALUES(?, ?, ?, ?)`,
		user.UserID, user.Username, user.Email, string(hash),
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			msg := mysqlErr.Message

			switch {
			case strings.Contains(msg, "uidx_username"):
				return logic.ErrUsernameTaken

			case strings.Contains(msg, "uidx_email"):
				return logic.ErrEmailTaken

			default:
				return logic.ErrDuplicateEntry
			}
		}
		return err
	}

	return nil
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
