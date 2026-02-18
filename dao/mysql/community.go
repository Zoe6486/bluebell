package mysql

import (
	"bluebell/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// // 这是一个构造函数，传入一个数据库连接 db，返回一个指向 CommunityDao 的指针。
// // 方便外部调用时用 NewCommunityDao(db) 来创建一个 CommunityDao 实例。
// func NewCommunityDao(db *sqlx.DB) *CommunityDao {
// 	return &CommunityDao{DB: db}
// }
// 如果像上面这样直接返回 *CommunityDao 指针， Service 层（业务层）代码就会长这样：
//	type CommunityService struct {
//	    Dao *mysql.CommunityDao // 强绑定！这里锁死了必须用 MySQL
//	}
// 这会产生一个大问题：单元测试。 大企业要求代码测试覆盖率。
// 如果 Service 强绑定了 MySQL DAO，你写测试时就必须真的启动一个 MySQL 数据库。
// 如果数据库挂了，或者别人改了数据，你的测试就失败了。

// interface 为了解耦和测试
// 1. 定义接口
// 这样 Service 层只需要知道有哪些方法，
// 只要你能提供这两个方法，你就是一个 CommunityStore。
// 不需要管底层是 MySQL， 还是 Redis， 还是 Mock 实现（测试用），还是内存实现
type CommunityStore interface {
	GetCommunityList() ([]*models.Community, error)
	GetCommunityDetailByID(id int64) (*models.CommunityDetail, error)
}

// 2. 结构体私有化（小写开头）
// 不让外部直接修改 DB 成员，强制通过构造函数实例化
// sqlx.DB 是 sqlx 库中对数据库连接的封装，基于 database/sql 做了增强，方便做数据库操作。
type communityDao struct {
	db *sqlx.DB
}

// 3. 构造函数返回接口类型
// 这是“大企业”常用的做法：返回接口而非具体指针
func NewCommunityDao(db *sqlx.DB) CommunityStore { // 返回的是接口,注意这里没加*，接口相当于类型信息 + 指向具体值的指针，所以不用加。具体类型 是communityDao 的指针，但藏在接口后面
	return &communityDao{db: db} // 塞进去的是具体的结构体指针
}

// --- 下面是具体实现 ---

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
			return nil, ErrorInvalidID
		}
		return nil, err
	}
	return &communityDetail, nil // 返回局部变量指针，会不会空？不会。Go 的逃逸分析会把它放到 heap。
}
