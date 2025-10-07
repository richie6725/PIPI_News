package aclCtrl

import (
	aclDaoModel "News/service/dao/daoModels/acl"
	aclMongoDao "News/service/dao/mongoDao/acl"
	aclRedisDao "News/service/dao/redisDao/acl"
	boAcl "News/service/internal/model/bo/acl"
	"News/service/internal/utils"
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"time"
)

type aclCtrl struct {
	pack aclCtrlPack
}

type aclCtrlPack struct {
	dig.In
	MongoNews *mongo.Database `name:"mongo_news"`
	RedisNews *redis.Client   `name:"redis_news"`
}

type AclCtrl interface {
	Get(ctx context.Context, args *boAcl.GetArgs) (*boAcl.GetReply, error)
	GetLogin(ctx context.Context, args *boAcl.GetArgs) (*boAcl.GetLoginReply, error)
	Update(ctx context.Context, args *boAcl.UpdateArgs) error
}

func NewAcl(pack aclCtrlPack) AclCtrl {
	return &aclCtrl{
		pack: pack,
	}
}

func (ctrl *aclCtrl) Get(ctx context.Context, args *boAcl.GetArgs) (*boAcl.GetReply, error) {
	aclDao := aclMongoDao.New(ctrl.pack.MongoNews)
	reply := &boAcl.GetReply{}

	user, err := aclDao.Get(ctx, args.User)
	if err != nil {
		return nil, err
	}

	if user != nil {
		reply.User = user
		return reply, nil
	}

	return nil, nil
}
func (ctrl *aclCtrl) GetLogin(ctx context.Context, args *boAcl.GetArgs) (*boAcl.GetLoginReply, error) {
	aclDao := aclMongoDao.New(ctrl.pack.MongoNews)
	aclRao := aclRedisDao.New(ctrl.pack.RedisNews)

	session, err := aclRao.Get(ctx, args.User.Username, args.User.Token)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return &boAcl.GetLoginReply{Session: session}, nil
	}

	user, err := aclDao.Get(ctx, args.User)
	if err != nil {
		return nil, err
	}

	if user != nil {
		token := utils.GenerateToken()

		session := &aclDaoModel.UserSession{
			Username: user.Username,
			Token:    token,
		}

		if err := aclRao.Set(ctx, session, time.Minute*30); err != nil {
			return nil, err		}

		return &boAcl.GetLoginReply{Session: session}, nil
	}

	return nil, nil
}

func (ctrl *aclCtrl) Update(ctx context.Context, args *boAcl.UpdateArgs) error {

	aclDao := aclMongoDao.New(ctrl.pack.MongoNews)

	err := aclDao.Update(ctx, aclDaoModel.Query{args.Query.BulkUserArgs, args.Query.CreatedAt})
	if err != nil {
		return err
	}

	return nil
}
