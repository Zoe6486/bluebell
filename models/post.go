package models

import "time"

// 内存对齐概念

// type Post struct {
// 	ID          int64     `json:"id,string" db:"post_id"`                            // 帖子id
// 	AuthorID    int64     `json:"author_id" db:"author_id"`                          // 作者id
// 	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"` // 社区id
// 	Status      int32     `json:"status" db:"status"`                                // 帖子状态
// 	Title       string    `json:"title" db:"title" binding:"required"`               // 帖子标题
// 	Content     string    `json:"content" db:"content" binding:"required"`           // 帖子内容
// 	CreateTime  time.Time `json:"create_time" db:"create_time"`                      // 帖子创建时间
// }
// CommunityID int64 json:"community_id" db:"community_id" binding:"required"``
// 优点：省事。一个结构体既能接收 Gin 的参数验证，又能直接丢给数据库存储。
// 缺点（致命伤）：权限污染。
// binding:"required" 是给 Gin 用的。万一你在某个业务逻辑里只想更新标题，不想传 ID，Gin 就会报错。
// json:"community_id" 是给前端用的。万一这个字段是内部逻辑计算出来的，不需要前端传，但因为有了这个 Tag，前端乱传一个值可能会覆盖你的业务逻辑。
//
// 零值陷阱：int64 的默认值是 0。如果 CommunityID 是必填的，且合法 ID 都是大于 0 的，那么 binding:"required" 会非常有用。
// 如果 0 也是个合法 ID，那你就得改用 *int64（指针）来区分“没传”和“传了0”。
// ???
//
// 在企业级开发中，我们通常会定义两个（甚至三个）结构体。
// 数据库模型（Model/Entity）Post
// 参数接收（Param/DTO）CreatePostParam
type Post struct {
	ID          int64     `db:"id"`
	PostID      int64     `db:"post_id"`
	Title       string    `db:"title"`
	Content     string    `db:"content"`
	AuthorID    int64     `db:"author_id"`
	CommunityID int64     `db:"community_id"`
	Status      int8      `db:"status"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	LikeCount   int64     `db:"like_count"` // 聚合计算字段，非表列
}

// PostStatus constants
const (
	PostStatusNormal  int8 = 1
	PostStatusDeleted int8 = 2
	PostStatusPending int8 = 3
)

// PostDetail is used for API responses, embedding enriched data
type PostDetail struct {
	*Post                // 嵌入 Post 全部字段
	AuthorName    string `db:"author_name"`    // users 表 JOIN 出来的???username?
	CommunityName string `db:"community_name"` // communities 表 JOIN
	// LikeCount     int64  `db:"like_count"` //Post写了就不写了
	DislikeCount int64 `db:"dislike_count"`
}

// CreatePostParams holds the data needed to create a new post
type CreatePostParams struct {
	PostID      int64  `json:"post_id"`
	Title       string `json:"title" binding:"required,max=128"`
	Content     string `json:"content" binding:"required"`
	AuthorID    int64  `json:"author_id"`
	CommunityID int64  `json:"community_id" binding:"required"`
}

// PostListParams holds pagination and filter params
type PostListParams struct {
	Page        int64  `json:"page" form:"page"`
	PageSize    int64  `json:"page_size" form:"page_size"`
	CommunityID int64  `json:"community_id" form:"community_id"`
	OrderBy     string `json:"order_by" form:"order_by"` // "time" | "score"
}
