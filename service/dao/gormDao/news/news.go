package news

import (
	newsDaoModel "News/service/dao/daoModels/news"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewsDao interface {
	Create(ctx context.Context, query []*newsDaoModel.News, sourceNews string) error
	Get(ctx context.Context, query []*newsDaoModel.News, sourceNews []string) ([]*newsDaoModel.News, error)
}

func New(db *gorm.DB) NewsDao {
	return &newsDao{
		db: db,
	}
}

type newsDao struct {
	db *gorm.DB
}

func (dao *newsDao) Create(ctx context.Context, query []*newsDaoModel.News, sourceNews string) error {

	err := dao.db.WithContext(ctx).Table(sourceNews + "_" + newsDaoModel.TableNews).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(&query).Error
	if err != nil {
		return fmt.Errorf("failed to insert news: %w", err)
	}

	return nil
}

func (dao *newsDao) Get(ctx context.Context, query []*newsDaoModel.News, sourceNews []string) ([]*newsDaoModel.News, error) {

	var (
		results []*newsDaoModel.News
	)

	for i, src := range sourceNews {
		var output newsDaoModel.News

		if err := dao.db.WithContext(ctx).Table(src+"_"+newsDaoModel.TableNews).
			Where(newsDaoModel.ID.AddField(), query[i].ID).Find(&output).Error; err != nil {
			return nil, fmt.Errorf("failed to get news from table: %w", err)
		}
		results = append(results, &output)
	}

	return results, nil
}
