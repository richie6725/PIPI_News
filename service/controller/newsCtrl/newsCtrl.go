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
	GetFilterOuterLayer(ctx context.Context, args *boNews.GetFilterOuterLayerArgs) (*boNews.GetFilterOuterLayerReply, error)
}

func NewNews(pack newsCtrlPack) NewsCtrl {
	return &newsCtrl{
		pack: pack,
	}
}

func (ctrl *newsCtrl) CreateOuterLayer(ctx context.Context, args *boNews.CreateOuterLayerArgs) error {

	newsDao := newsOuterLayer.New(ctrl.pack.Postgres)

	err := newsDao.Create(ctx, args.Query, args.SourceNews)
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *newsCtrl) GetOuterLayer(ctx context.Context, args *boNews.GetOuterLayerArgs) (*boNews.GetOuterLayerReply, error) {
	newsDao := newsOuterLayer.New(ctrl.pack.Postgres)
	outerLayers, err := newsDao.Get(ctx, args.Query)

	if err != nil {
		return nil, err
	}

	result := &boNews.GetOuterLayerReply{
		Data: outerLayers,
	}

	return result, nil
}

func (ctrl *newsCtrl) GetFilterOuterLayer(ctx context.Context, args *boNews.GetFilterOuterLayerArgs) (*boNews.GetFilterOuterLayerReply, error) {
	newsDao := newsOuterLayer.New(ctrl.pack.Postgres)
	filteredOuters, err := newsDao.GetFilter(ctx, args.Query, args.Pagination, args.TimeInterval)

	if err != nil {
		return nil, err
	}

	reply := &boNews.GetFilterOuterLayerReply{
		Data:       filteredOuters,
		Pagination: args.Pagination,
	}

	return reply, nil
}
