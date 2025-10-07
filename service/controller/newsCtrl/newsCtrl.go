package newsCtrl

import (
	"News/service/dao/gormDao/newsOuterLayer"
	boNews "News/service/internal/model/bo/news"
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type newsCtrl struct {
	pack newsCtrlPack
}

type newsCtrlPack struct {
	dig.In
	Postgres *gorm.DB `name:"postgres_news"`
}

type NewsCtrl interface {
	CreateOuterLayer(ctx context.Context, args *boNews.CreateOuterLayerArgs) error
	GetOuterLayer(ctx context.Context, args *boNews.GetOuterLayerArgs) (*boNews.GetOuterLayerReply, error)
}

func NewNews(pack newsCtrlPack) NewsCtrl {
	return &newsCtrl{
		pack: pack,
	}
}

func (ctrl *newsCtrl) CreateOuterLayer(ctx context.Context, args *boNews.CreateOuterLayerArgs) error {

	aclDao := newsOuterLayer.New(ctrl.pack.Postgres)

	err := aclDao.Create(ctx, args.Query, args.SourceNews)
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *newsCtrl) GetOuterLayer(ctx context.Context, args *boNews.GetOuterLayerArgs) (*boNews.GetOuterLayerReply, error) {
	aclDao := newsOuterLayer.New(ctrl.pack.Postgres)
	results, err := aclDao.Get(ctx, args.Query)

	if err != nil {
		return nil, err
	}

	reply := &boNews.GetOuterLayerReply{
		Data: results,
	}

	return reply, nil
}
