package database

import "time"

type TimeInterval struct {
	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`
}

type Pagenation struct {
	PageSize   int `bson:"page_size" json:"page_size"`
	PageNumber int `bson:"page_number" json:"page_number"`
}
