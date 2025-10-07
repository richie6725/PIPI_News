package news

import newsDaoModel "News/service/dao/daoModels/news"

type CreateOuterLayerArgs struct {
	Query      []*newsDaoModel.OuterLayer
	SourceNews string
}

type CreateOuterLayerReply struct{}

type GetOuterLayerArgs struct {
	Query []*newsDaoModel.OuterLayer
}

type GetOuterLayerReply struct {
	Data []*newsDaoModel.OuterLayer
}

type CreateContentLayerArgs struct {
	Query []*newsDaoModel.ContentLayer
}

type CreateContentLayerReply struct {
}
