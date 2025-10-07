package newsOuterLayer

import (
	newsDaoModel "News/service/dao/daoModels/news"
	daoModel "News/service/internal/database"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type NewsOuterLayerDao interface {
	Get(ctx context.Context, query []*newsDaoModel.OuterLayer) ([]*newsDaoModel.OuterLayer, error)
	Create(ctx context.Context, query []*newsDaoModel.OuterLayer, source string) error
	GetFilter(ctx context.Context, query *newsDaoModel.OuterLayer, pagination daoModel.Pagenation, timeInterval daoModel.TimeInterval) ([]*newsDaoModel.OuterLayer, error)
}

func New(db *gorm.DB) NewsOuterLayerDao {
	return &newsOuterLayerDao{
		db: db,
	}
}

type newsOuterLayerDao struct {
	db *gorm.DB
}

func (dao *newsOuterLayerDao) Get(ctx context.Context, query []*newsDaoModel.OuterLayer) ([]*newsDaoModel.OuterLayer, error) {
	if len(query) == 0 {
		return []*newsDaoModel.OuterLayer{}, fmt.Errorf("no data to insert: query slice is empty")
	}

	db := dao.db.WithContext(ctx).Table(newsDaoModel.Table)

	for i, item := range query {
		condition := dao.db.Where("source_news = ?", item.SourceNews).
			Where("news_id = ?", item.NewsID).
			Where("title = ?", item.Title)
		if i == 0 {
			db = db.Where(condition)
		} else {
			db = db.Or(condition)
		}
	}

	var results []*newsDaoModel.OuterLayer
	if err := db.Order("news_id DESC").Find(&results).Error; err != nil {
		return nil, fmt.Errorf("database error on bulk get: %w", err)
	}

	return results, nil
}

func (dao *newsOuterLayerDao) Create(ctx context.Context, query []*newsDaoModel.OuterLayer, source string) error {
	if len(query) == 0 {
		return fmt.Errorf("no data to insert: query slice is empty")
	}

	for _, item := range query {
		item.SourceNews = source
	}

	if err := dao.db.WithContext(ctx).
		Table(newsDaoModel.Table).
		Create(&query).Error; err != nil {
		return fmt.Errorf("failed to insert news outer layer: %w", err)
	}

	return nil
}

func (dao *newsOuterLayerDao) GetFilter(ctx context.Context, query *newsDaoModel.OuterLayer, pagination daoModel.Pagenation, timeInterval daoModel.TimeInterval) ([]*newsDaoModel.OuterLayer, error) {

	//if err:=dao.db.WithContext(ctx).

	return nil, nil
}
