package mysql

import (
	"bluebell/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// 实现接口的方法
// 结构体私有化（小写开头）
// 不让外部直接修改 DB 成员，强制通过构造函数实例化
// sqlx.DB 是 sqlx 库中对数据库连接的封装，基于 database/sql 做了增强，方便做数据库操作。
type communityDao struct {
	db *sqlx.DB
}

// 构造函数：虽然它返回 *communityDao，但因为它实现了 Logic 里的接口，
// 所以可以被直接注入到 Logic 中。
func NewCommunityDao(db *sqlx.DB) *communityDao {
	return &communityDao{db: db}
}

func (dao *communityDao) GetCommunityList() ([]*models.Community, error) {
	var communityList []*models.Community
	sqlStr := `select community_id, community_name from community`
	// 注意区分 Get（查单条），有ErrNoRows
	// Select 需要传指针
	// Select 把“没数据”已经用空切片表达，没有ErrNoRows这个错误
	if err := dao.db.Select(&communityList, sqlStr); err != nil {
		return nil, err
	}
	return communityList, nil // 注意前面那个加&，这个不加&,再加类型就变成*[]*models.Community了，和返回值要求的不一样了
}

func (dao *communityDao) GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	var communityDetail models.CommunityDetail
	sqlStr := `select community_id, community_name, introduction, create_time from community where community_id = ?`
	err := dao.db.Get(&communityDetail, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// 	return nil, ErrorInvalidID
			return nil, err //到底在哪里定义错误？？？
		}
		return nil, err
	}
	return &communityDetail, nil // 返回局部变量指针，会不会空？不会。Go 的逃逸分析会把它放到 heap。
}
