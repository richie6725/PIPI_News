package news

import (
	newsDaoModel "News/service/dao/daoModels/news"
	"News/service/internal/database"
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewsDao interface {
	Create(ctx context.Context, query []*newsDaoModel.News, sourceNews string) error
	Get(ctx context.Context, query []*newsDaoModel.News, sourceNews []string) ([]*newsDaoModel.News, error)
	GetFilter(ctx context.Context, query *newsDaoModel.News, sourceNews string, timeInterval *database.TimeInterval, pageNation *database.Pagenation) ([]*newsDaoModel.News, error)
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

//func (dao *newsDao) GetFilter(ctx context.Context, query *newsDaoModel.News, sourceNews string, timeInterval *database.TimeInterval, pageNation *database.Pagenation) ([]*newsDaoModel.News, error) {
//
//	var (
//		results []*newsDaoModel.News
//	)
//
//	if err := dao.db.WithContext(ctx).Table(sourceNews+"_"+newsDaoModel.TableNews).
//		Where(newsDaoModel.Category.AddField(), query.Category).
//		Where(newsDaoModel.TitleTags.OverlapsField(), query.TitleTags).
//		Where(newsDaoModel.ReleaseTime.MoreThanEqualField(), timeInterval.StartTime).
//		Where(newsDaoModel.ReleaseTime.LessThanEqualField(), timeInterval.EndTime).Find(&results).Error; err != nil {
//		return nil, fmt.Errorf("failed to get news from table: %w", err)
//	}
//
//	return results, nil
//}

func (dao *newsDao) GetFilter(
	ctx context.Context,
	query *newsDaoModel.News,
	sourceNews string,
	timeInterval *database.TimeInterval,
	pageNation *database.Pagenation,
) ([]*newsDaoModel.News, error) {

	var (
		results []*newsDaoModel.News
	)

	db := dao.db.WithContext(ctx).Table(sourceNews + "_" + newsDaoModel.TableNews)

	// -------------------------------
	// 動態條件組合
	// -------------------------------
	filterCount := 0

	if query.Category != "" {
		db = db.Where(newsDaoModel.Category.AddField(), query.Category)
		filterCount++
	}

	if len(query.TitleTags) > 0 {
		db = db.Where(newsDaoModel.TitleTags.OverlapsField(), query.TitleTags)
		filterCount++
	}

	if !timeInterval.StartTime.IsZero() {
		db = db.Where(newsDaoModel.ReleaseTime.MoreThanEqualField(), timeInterval.StartTime)
		filterCount++
	}

	if !timeInterval.EndTime.IsZero() {
		db = db.Where(newsDaoModel.ReleaseTime.LessThanEqualField(), timeInterval.EndTime)
		filterCount++
	}

	// -------------------------------
	// 安全檢查：至少要有一個篩選條件
	// -------------------------------
	if filterCount == 0 {
		return nil, fmt.Errorf("no valid filter provided — at least one condition is required")
	}

	// -------------------------------
	// 分頁處理
	// -------------------------------
	if pageNation.PageSize > 0 {
		offset := (pageNation.Page - 1) * pageNation.PageSize
		db = db.Offset(offset).Limit(pageNation.PageSize)
	}

	// -------------------------------
	// 查詢執行
	// -------------------------------
	if err := db.Order(newsDaoModel.ReleaseTime.String() + " DESC").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get news from table: %w", err)
	}

	return results, nil
}
