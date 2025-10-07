package acl

import (
	"News/service/dao/daoModels/acl"
	"News/service/dao/mongoDao"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	aclCollection = "acl"
)

type AclDao interface {
	Get(ctx context.Context, user aclDaoModel.User) (*aclDaoModel.User, error)
	Update(ctx context.Context, query aclDaoModel.Query) error
}

func New(db *mongo.Database) AclDao {
	dao := &aclDao{
		collection: db.Collection(aclCollection),
	}
	return dao
}

type aclDao struct {
	collection *mongo.Collection
}

func (dao *aclDao) Get(ctx context.Context, user aclDaoModel.User) (*aclDaoModel.User, error) {
	pipe := mongoDao.NewStageBuilder().
		AddMatch(buildMatchQueries(user)).
		Generate()

	cursor, err := dao.collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result aclDaoModel.User
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, err

}

//寫入單筆
//func (dao *aclDao) Update(ctx context.Context, model aclDaoModel.User) error {
//
//	filter := mongoDao.NewMatchBuilder().
//		AddEqual(aclDaoModel.Username, model.Username).Generate()
//
//	if len(filter) == 0 {
//		return fmt.Errorf("missing pk, model: %+v", model)
//	}
//	doc := bson.M{"$set": model}
//
//	output, err := dao.collection.UpdateOne(ctx, filter, doc, options.Update().SetUpsert(true))
//	fmt.Printf("MatchedCount: %d, ModifiedCount: %d, UpsertedID: %v\n",
//		output.MatchedCount,
//		output.ModifiedCount,
//		output.UpsertedID,
//	)
//	return err
//}

func (dao *aclDao) Update(ctx context.Context, query aclDaoModel.Query) error {

	writes := make([]mongo.WriteModel, len(query.BulkUserArgs))

	for i := range writes {
		filter := mongoDao.NewMatchBuilder().
			AddEqual(aclDaoModel.Username, query.BulkUserArgs[i].Username).Generate()
		if len(filter) == 0 {
			return fmt.Errorf("missing pk, model: %+v", query.BulkUserArgs[i])
		}
		doc := bson.M{"$set": aclDaoModel.UserWithMeta{
			BulkUserArg: query.BulkUserArgs[i],
			CreatedAt:   query.CreatedAt,
		}}
		writes[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(doc).SetUpsert(true)

	}

	_, err := dao.collection.BulkWrite(ctx, writes, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return err
	}
	return nil
}

func buildMatchQueries(user aclDaoModel.User) []bson.E {
	queries := mongoDao.NewMatchBuilder().
		AddEqual(aclDaoModel.Username, user.Username).
		AddEqual(aclDaoModel.Password, user.Password)

	return queries.Generate()
}
