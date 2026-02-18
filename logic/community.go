package logic

import (
	"bluebell/models"
	"database/sql"
	"errors"
)

// 💡 接口定义在 Logic 层！
// 它描述了 Logic 需要什么样的能力，不关心底层是 MySQL 还是文件
type CommunityStore interface {
	GetCommunityList() ([]*models.Community, error)
	GetCommunityDetailByID(id int64) (*models.CommunityDetail, error)
}

type CommunityLogic struct {
	// 注入接口，完全不依赖 mysql 包
	store CommunityStore
}

func NewCommunityLogic(s CommunityStore) *CommunityLogic {
	return &CommunityLogic{store: s}
}

func (l *CommunityLogic) GetCommunityList() ([]*models.Community, error) {
	return l.store.GetCommunityList()
}

var ErrCommunityNotFound = errors.New("community not found")

func (l *CommunityLogic) GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// return l.store.GetCommunityDetailByID(id)
	data, err := l.store.GetCommunityDetailByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCommunityNotFound
		}
		return nil, err
	}
	return data, nil
}
