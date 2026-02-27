package mysql

import "github.com/jmoiron/sqlx"

type postDao struct {
	db *sqlx.DB
}

func NewPostDao(db *sqlx.DB) *postDao {
	return &postDao{db: db}
}
