package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// CommunityLogic 业务逻辑结构体
type CommunityLogic struct {
	// 💡 重点：这里存的是 DAO 的接口，而不是具体的 MySQL 实现
	dao mysql.CommunityStore
}

// NewCommunityLogic 构造函数，把 DAO 注入进来
func NewCommunityLogic(dao mysql.CommunityStore) *CommunityLogic { //struct 是值类型,所以得返回指针避免返回副本
	return &CommunityLogic{dao: dao}
}

func (l *CommunityLogic) GetCommunityList() ([]*models.Community, error) {
	// 这里可以写业务逻辑，比如：权限校验、缓存读取等
	// 现在直接调用 DAO 获取数据
	return l.dao.GetCommunityList()
}
func (l *CommunityLogic) GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return l.dao.GetCommunityDetailByID(id)
}
