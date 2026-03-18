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

//内存对齐优化后，注释了，因为不想改了
// type User struct {
//     // 8 字节对齐的字段放前面（天然对齐）
//     ID        int64     `db:"id"`
//     UserID    int64     `db:"user_id"`

//     CreatedAt time.Time `db:"created_at"`  // 24 字节，8 对齐
//     UpdatedAt time.Time `db:"updated_at"`  // 24 字节，8 对齐

//     Username  string    `db:"username"`    // 16 字节，8 对齐
//     Email     string    `db:"email"`
//     Password  string    `db:"password"`
//     Token     string    `db:"-"`           // runtime only，也放这里没问题

//     // 小的放最后（甚至可以多个小字段挤在一起）
//     Status    int8      `db:"status"`
//     // 如果以后有其他 int8/bool/uint8，可以继续加在这里
// }

// // Token 分离出去
// type User struct {
// 	ID     int64 `db:"id"`
// 	UserID int64 `db:"user_id"`

// 	CreatedAt time.Time `db:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at"`

// 	Username string `db:"username"`
// 	Email    string `db:"email"`
// 	Password string `db:"password"`

// 	Status int8 `db:"status"`
// }
// type UserWithToken struct {
// 	User
// 	Token string `json:"token"`
// }
