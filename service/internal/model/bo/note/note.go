package note

import (
	noteDaoModel "News/service/dao/daoModels/note"
)

type GetArgs struct {
	Query noteDaoModel.Query
}

type GetReply struct {
	BulkPriceRecordArgs []*noteDaoModel.PriceRecord
}

type UpdateArgs struct {
	BulkPriceRecordArgs []*noteDaoModel.PriceRecord
	IsUpsert            bool
}

type UpdateReply struct {
}
