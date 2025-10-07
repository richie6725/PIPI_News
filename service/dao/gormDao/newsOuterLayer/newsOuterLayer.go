package newsOuterLayer

import (
	newsDaoModel "News/service/dao/daoModels/news"
	"context"
	"gorm.io/gorm"
)

type NewsOuterLayerDao interface {
	Get()
	Create(ctx context.Context, query []*newsDaoModel.OuterLayer, source string) error
}

func New(db *gorm.DB) NewsOuterLayerDao {
	return &newsOuterLayerDao{
		db: db,
	}
}

type newsOuterLayerDao struct {
	db *gorm.DB
}

func (dao *newsOuterLayerDao) Get() {

}

func (dao *newsOuterLayerDao) Create(ctx context.Context, query []*newsDaoModel.OuterLayer, source string) error {
	if len(query) == 0 {
		return nil // 如果沒有資料，直接返回
	}

	// 在將資料傳遞給 GORM 之前，為每一筆資料設定 `SourceNews` 欄位
	for _, item := range query {
		item.SourceNews = source
	}

	return dao.db.WithContext(ctx).
		Table(newsDaoModel.Table).
		Create(&query).
		Error
}
