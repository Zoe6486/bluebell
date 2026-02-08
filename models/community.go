package models

import "time"

// 一般用于列表接口，只需要返回最核心、最精简的信息。
// 避免某个接口暴露过多不需要的数据。
// 比如社区的ID和名字，用于展示列表，数据量小，网络传输快，性能更好。
type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

// ommunityDetail：用于详情接口，返回更多信息。
// 比如用户查看单个社区的详细信息时用。
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
