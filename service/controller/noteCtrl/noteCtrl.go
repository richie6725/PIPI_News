package noteCtrl

import (
	noteMongoDao "News/service/dao/mongoDao/note"
	boNote "News/service/internal/model/bo/note"
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
)

type noteCtrl struct {
	pack noteCtrlPack
}

type noteCtrlPack struct {
	dig.In
	MongoNews *mongo.Database `name:"mongo_news"`
	RedisNews *redis.Client   `name:"redis_news"`
}

type NoteCtrl interface {
	Update(ctx context.Context, args *boNote.UpdateArgs) error
	Get(ctx context.Context, args *boNote.GetArgs) (*boNote.GetReply, error)
}

func NewNote(pack noteCtrlPack) NoteCtrl {
	return &noteCtrl{
		pack: pack,
	}
}

func (ctrl *noteCtrl) Update(ctx context.Context, args *boNote.UpdateArgs) error {
	noteDao := noteMongoDao.New(ctrl.pack.MongoNews)

	err := noteDao.Update(ctx, args.BulkPriceRecordArgs, args.IsUpsert)
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *noteCtrl) Get(ctx context.Context, args *boNote.GetArgs) (*boNote.GetReply, error) {
	noteDao := noteMongoDao.New(ctrl.pack.MongoNews)

	note, err := noteDao.Get(ctx, args.Query)

	if err != nil {
		return nil, err
	}

	reply := &boNote.GetReply{
		BulkPriceRecordArgs: note,
	}

	return reply, nil

}
