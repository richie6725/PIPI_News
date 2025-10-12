package newsOuterLayer

import (
	newsDaoModel "News/service/dao/daoModels/news"
	daoModel "News/service/internal/database"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		condition := dao.db.Where(newsDaoModel.SourceNews.AddField(), item.SourceNews).
			Where(newsDaoModel.ID.AddField(), item.NewsID).
			Where(newsDaoModel.Title.AddField(), item.Title)
		if i == 0 {
			db = db.Where(condition)
		} else {
			db = db.Or(condition)
		}
	}

	var results []*newsDaoModel.OuterLayer
	if err := db.Order(newsDaoModel.ID + "DESC").Find(&results).Error; err != nil {
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
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "source_news"}, {Name: "news_id"}},
			UpdateAll: true,
		}).
		Create(&query); err != nil {
		return fmt.Errorf("failed to insert news outer layer: %w", err)
	}

	return nil
}
func (dao *newsOuterLayerDao) GetFilter(ctx context.Context, query *newsDaoModel.OuterLayer, pagination daoModel.Pagenation, timeInterval daoModel.TimeInterval) ([]*newsDaoModel.OuterLayer, error) {

	db := dao.db.WithContext(ctx).Table(newsDaoModel.Table)

	var hasConditions bool

	if query.SourceNews != "" {
		db = db.Where(newsDaoModel.SourceNews.AddField(), query.SourceNews)
		hasConditions = true
	}

	if query.Category != "" {
		db = db.Where(newsDaoModel.Category.AddField(), query.Category)
		hasConditions = true
	}

	if len(query.Tags) > 0 {
		db = db.Where(newsDaoModel.Tags.ContainField(), query.Tags)
		hasConditions = true
	}

	if !timeInterval.StartTime.IsZero() {
		db = db.Where(newsDaoModel.ReleaseTime.MoreThanEqualField(), timeInterval.StartTime)
		hasConditions = true
	}
	if !timeInterval.EndTime.IsZero() {
		db = db.Where(newsDaoModel.ReleaseTime.LessThanEqualField(), timeInterval.EndTime)
		hasConditions = true
	}

	if !hasConditions {
		return nil, errors.New("at least one filter condition is required")
	}

	pageSkip := pagination.PageSize * (pagination.Page - 1)

	if pagination.PageSize > 0 {
		db = db.Offset(pageSkip).Limit(pagination.PageSize)
	}

	var results []*newsDaoModel.OuterLayer
	if err := db.Order(newsDaoModel.ID + "DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
