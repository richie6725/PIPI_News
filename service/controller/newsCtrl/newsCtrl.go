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
}

func NewNews(pack newsCtrlPack) NewsCtrl {
	return &newsCtrl{
		pack: pack,
	}
}

func (ctrl *newsCtrl) CreateOuterLayer(ctx context.Context, args *boNews.CreateOuterLayerArgs) error {

	aclDao := newsOuterLayer.New(ctrl.pack.Postgres)

	err := aclDao.Create(ctx, args.Query, args.Source)
	if err != nil {
		return err
	}

	return nil
}
