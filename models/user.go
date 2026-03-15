// package models

//	type User struct {
//		UserID   int64  `db:"user_id"`
//		Username string `db:"username"`
//		Password string `db:"password"`
//		Token    string
//	}
package models

import "time"

type User struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	// Runtime only — never persisted
	Token string `db:"-"`
}

const (
	UserStatusActive    int8 = 1
	UserStatusSuspended int8 = 2
	UserStatusDeleted   int8 = 3
)
