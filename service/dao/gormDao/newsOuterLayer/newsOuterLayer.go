package newsOuterLayer

import (
	newsDaoModel "News/service/dao/daoModels/news"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type NewsOuterLayerDao interface {
	Get(ctx context.Context, query []*newsDaoModel.OuterLayer) ([]*newsDaoModel.OuterLayer, error)
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

func (dao *newsOuterLayerDao) Get(ctx context.Context, query []*newsDaoModel.OuterLayer) ([]*newsDaoModel.OuterLayer, error) {
	if len(query) == 0 {
		return []*newsDaoModel.OuterLayer{}, nil
	}

	// --- 步驟 1: 建立一個複雜的 OR 查詢 ---
	db := dao.db.WithContext(ctx).Table(newsDaoModel.Table)

	// 將多個 (source_news = ? AND news_id = ? AND title = ?) 條件用 OR 連接
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

	// --- 步驟 2: 執行單次查詢，獲取所有匹配的結果 ---
	var dbResults []*newsDaoModel.OuterLayer
	if err := db.Find(&dbResults).Error; err != nil {
		return nil, fmt.Errorf("database error on bulk get: %w", err)
	}

	// --- 步驟 3: 將查詢結果放入一個 map 中，方便快速查找 ---
	// 我們需要一個唯一的 key 來標識每一筆資料
	resultsMap := make(map[string]*newsDaoModel.OuterLayer)
	for _, result := range dbResults {
		// 使用 "SourceNews-NewsID-Title" 作為 map 的 key
		key := fmt.Sprintf("%s-%d-%s", result.SourceNews, result.NewsID, result.Title)
		resultsMap[key] = result
	}

	// --- 步驟 4: 遍歷原始的 query slice，從 map 中取出結果，組裝成最終的、有序的 slice ---
	finalResults := make([]*newsDaoModel.OuterLayer, len(query))
	for i, item := range query {
		key := fmt.Sprintf("%s-%d-%s", item.SourceNews, item.NewsID, item.Title)

		// 從 map 中查找對應的結果
		if foundItem, ok := resultsMap[key]; ok {
			// 如果找到了，放入最終結果的對應位置
			finalResults[i] = foundItem
		} else {
			// 如果沒找到，對應位置會是預設的 nil，符合需求
			finalResults[i] = nil
		}
	}

	return finalResults, nil
}

func (dao *newsOuterLayerDao) Create(ctx context.Context, query []*newsDaoModel.OuterLayer, source string) error {
	if len(query) == 0 {
		return nil
	}

	for _, item := range query {
		item.SourceNews = source
	}

	return dao.db.WithContext(ctx).
		Table(newsDaoModel.Table).
		Create(&query).
		Error
}
