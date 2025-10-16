package newsCtrl

import (
	"News/service/dao/gormDao/news"
	boNews "News/service/internal/model/bo/news"
	"News/service/internal/utils"
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
	CreateNews(ctx context.Context, args *boNews.CreateNewsArgs) error
	GetNews(ctx context.Context, args *boNews.GetNewsArgs) (*boNews.GetNewsReply, error)
	GetNewsFilter(ctx context.Context, args *boNews.GetNewsFilterArgs) (*boNews.GetNewsFilterReply, error)
}

func NewNews(pack newsCtrlPack) NewsCtrl {
	return &newsCtrl{
		pack: pack,
	}
}

func (ctrl *newsCtrl) CreateNews(ctx context.Context, args *boNews.CreateNewsArgs) error {
	newsDao := news.New(ctrl.pack.Postgres)

	err := newsDao.Create(ctx, args.Query, args.SourceNews)
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *newsCtrl) GetNews(ctx context.Context, args *boNews.GetNewsArgs) (*boNews.GetNewsReply, error) {

	newsDao := news.New(ctrl.pack.Postgres)

	eachNewsData, err := newsDao.Get(ctx, args.Query, args.SourceNews)

	if err != nil {
		return nil, err
	}

	reply := &boNews.GetNewsReply{
		Data: eachNewsData,
	}

	return reply, nil
}

func (ctrl *newsCtrl) GetNewsFilter(ctx context.Context, args *boNews.GetNewsFilterArgs) (*boNews.GetNewsFilterReply, error) {

	newsDao := news.New(ctrl.pack.Postgres)
	filteredData, err := newsDao.GetFilter(ctx, args.Query, args.SourceNews, args.TimeInterval, args.Pagination)

	if err != nil {
		return nil, err
	}

	reply := &boNews.GetNewsFilterReply{
		Data:       filteredData,
		Pagination: utils.BuildPagination(args.Pagination.Page, args.Pagination.PageSize, len(filteredData)),
	}

	return reply, nil
}
