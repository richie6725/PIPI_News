package news

import newsDaoModel "News/service/dao/daoModels/news"

type CreateOuterLayerArgs struct {
	Query  []*newsDaoModel.OuterLayer
	Source string
}

type CreateOuterLayerReply struct{}

type CreateContentLayerArgs struct {
	Query []*newsDaoModel.ContentLayer
}

type CreateContentLayerReply struct {
}
