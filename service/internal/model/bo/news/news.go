package news

import (
	newsDaoModel "News/service/dao/daoModels/news"
	daoModel "News/service/internal/database"
	"News/service/internal/utils"
)

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

type GetNewsFilterArgs struct {
	Query        *newsDaoModel.News
	SourceNews   string
	Pagination   *daoModel.Pagenation
	TimeInterval *daoModel.TimeInterval
}

type GetNewsFilterReply struct {
	Data       []*newsDaoModel.News
	Pagination utils.Pagination
}

type CreateContentLayerArgs struct {
	Query []*newsDaoModel.ContentLayer
}

type CreateContentLayerReply struct {
}
