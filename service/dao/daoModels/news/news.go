package news

import (
	"github.com/lib/pq"
	"time"
)

const (
	Table = "news_outer_layer"
)

type BaseNews struct {
	SourceNews  string         `json:"-"`
	Title       string         `json:"title"`
	DetailURL   string         `json:"detail_url"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	ImgURL      string         `json:"img_url"`
	NewsID      int            `json:"news_id"`
	ReleaseTime time.Time      `json:"release_time"`
	Category    string         `json:"category"`
	CategoryURL string         `json:"category_url"`
}

type ContentLayer struct {
	BaseNews
	ImgString string `json:"img_string"`
	Reporter  string `json:"reporter"`
	Content   string `json:"content"`
}

type OuterLayer struct {
	BaseNews
}
