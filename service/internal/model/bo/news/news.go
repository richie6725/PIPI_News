package news

import (
	newsDaoModel "News/service/dao/daoModels/news"
	daoModel "News/service/internal/database"
)

type CreateOuterLayerArgs struct {
	Query      []*newsDaoModel.OuterLayer
	SourceNews string
}

type CreateOuterLayerReply struct{}

type CreateNewsArgs struct {
	Query      []*newsDaoModel.News
	SourceNews string
}

type CreateNewsReply struct {
}

type GetNewsArgs struct {
	Query      []*newsDaoModel.News
	SourceNews []string
}

type GetNewsReply struct {
	Data []*newsDaoModel.News
}

type GetOuterLayerArgs struct {
	Query []*newsDaoModel.OuterLayer
}

type GetOuterLayerReply struct {
	Data []*newsDaoModel.OuterLayer
}

type GetFilterOuterLayerArgs struct {
	Query        *newsDaoModel.OuterLayer
	Pagination   daoModel.Pagenation
	TimeInterval daoModel.TimeInterval
}

type GetFilterOuterLayerReply struct {
	Data       []*newsDaoModel.OuterLayer
	Pagination daoModel.Pagenation
}

type CreateContentLayerArgs struct {
	Query []*newsDaoModel.ContentLayer
}

type CreateContentLayerReply struct {
}
