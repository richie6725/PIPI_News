package note

import (
	noteDaoModel "News/service/dao/daoModels/note"
	"News/service/dao/mongoDao"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const noteCollection = "note"

type NoteDao interface {
	Update(ctx context.Context, models []*noteDaoModel.PriceRecord, isUpsert bool) error
	Get(ctx context.Context, query noteDaoModel.Query) ([]*noteDaoModel.PriceRecord, error)
}

type noteDao struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) NoteDao {
	dao := &noteDao{
		collection: db.Collection(noteCollection),
	}
	return dao
}

func (dao *noteDao) Update(ctx context.Context, models []*noteDaoModel.PriceRecord, isUpsert bool) error {
	writes := make([]mongo.WriteModel, len(models))
	for i := range writes {
		filter := mongoDao.NewMatchBuilder().
			AddEqual(noteDaoModel.FieldMarketName, models[i].Market.Name).
			AddEqual(noteDaoModel.FieldProductID, models[i].Product.ProductID).Generate()
		if len(filter) == 0 {
			return fmt.Errorf("invalid filter: market=%s productID=%s", models[i].Market.Name, models[i].Product.ProductID)
		}
		doc := bson.M{"$set": models[i]}
		update := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(doc).SetUpsert(isUpsert)
		writes[i] = update
	}
	_, err := dao.collection.BulkWrite(ctx, writes, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return err
	}

	return nil
}

func (dao *noteDao) Get(ctx context.Context, query noteDaoModel.Query) ([]*noteDaoModel.PriceRecord, error) {

	pipe := mongoDao.NewStageBuilder().
		AddMatch(buildMatchQueries(query)).Generate()

	cursor, err := dao.collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*noteDaoModel.PriceRecord
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode cursor results: %w", err)
	}

	return results, nil
}

func buildMatchQueries(query noteDaoModel.Query) bson.D {

	queries := mongoDao.NewMatchBuilder().
		AddOr(buildMarketProductMatch(query.BulkPriceRecordArgs)).
		AddBetween(noteDaoModel.FieldCreatedAt, query.TimeInterval.StartTime, query.TimeInterval.EndTime).
		Generate()

	return queries

}

func buildMarketProductMatch(bulkPriceRecordArgs []noteDaoModel.PriceRecord) bson.A {
	var orConditions bson.A
	for _, record := range bulkPriceRecordArgs {
		pairCondition := bson.M{
			noteDaoModel.FieldMarketID.String():  record.Market.MarketID,
			noteDaoModel.FieldProductID.String(): record.Product.ProductID,
		}
		orConditions = append(orConditions, pairCondition)
	}
	return orConditions
}

func buildMatchQueries_old(query noteDaoModel.Query) bson.D {
	// 初始化一個空的篩選器

	var filter bson.D

	// --- Part 1: 建立 MarketID 和 ProductID 的組合查詢 ---
	// 只有當 BulkPriceRecordArgs 列表不為空時，才需要加入 $or 條件
	if len(query.BulkPriceRecordArgs) > 0 {
		// $or 的值是一個陣列，陣列中每個元素都是一個查詢文件
		// e.g., "$or": [ { "market.market_id": "M1", "product.product_id": "P1" }, { ... } ]
		var orConditions []bson.D

		for _, record := range query.BulkPriceRecordArgs {
			// 為每一對 MarketID 和 ProductID 建立一個獨立的查詢條件
			// 這個條件要求 market_id 和 product_id 必須同時匹配
			pairCondition := bson.D{
				{Key: "market.market_id", Value: record.Market.MarketID},
				{Key: "product.product_id", Value: record.Product.ProductID},
			}
			orConditions = append(orConditions, pairCondition)
		}

		// 將 $or 條件加入到主篩選器中
		filter = append(filter, bson.E{Key: "$or", Value: orConditions})
	}

	// --- Part 2: 建立時間區間查詢 ---
	// 建立一個子文件來存放 $gte 和 $lte
	// e.g., "meta.created_at": { "$gte": startTime, "$lte": endTime }
	timeConditions := bson.D{}

	// 如果有提供開始時間，則加入 $gte (greater than or equal) 條件
	if !query.TimeInterval.StartTime.IsZero() {
		timeConditions = append(timeConditions, bson.E{Key: "$gte", Value: query.TimeInterval.StartTime})
	}

	// 如果有提供結束時間，則加入 $lte (less than or equal) 條件
	if !query.TimeInterval.EndTime.IsZero() {
		timeConditions = append(timeConditions, bson.E{Key: "$lte", Value: query.TimeInterval.EndTime})
	}

	// 只有當時間條件存在時，才將 meta.created_at 篩選器加入到主篩選器
	if len(timeConditions) > 0 {
		filter = append(filter, bson.E{Key: "meta.created_at", Value: timeConditions})
	}

	return filter
}
