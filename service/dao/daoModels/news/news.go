package newsDaoModel

import (
	"News/service/dao/gormDao"
	"github.com/lib/pq"
	"time"
)

const (
	Table     = "news_outer_layer"
	TableNews = "news"
)

type FieldName = gormDao.FieldName

const (
	SourceNews  FieldName = "source_news"
	ID          FieldName = "id"
	Title       FieldName = "title"
	DetailURL   FieldName = "detail_url"
	Tags        FieldName = "tags"
	ImgURL      FieldName = "img_url"
	ReleaseTime FieldName = "release_time"
	Category    FieldName = "category"
	CategoryURL FieldName = "category_url"
)

type BaseNews struct {
	SourceNews  string         `json:"source_news"`
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

type News struct {
	ID             int    `gorm:"column:id;primaryKey" json:"id"`
	Title          string `gorm:"column:title" json:"title"`
	DetailURL      string `gorm:"column:detail_url" json:"detail_url"`
	OutPageImgURL  string `gorm:"column:out_page_img_url" json:"out_page_img_url"`
	Category       string `gorm:"column:category" json:"category"`
	CategoryURL    string `gorm:"column:category_url" json:"category_url"`
	HeaderImgURL   string `gorm:"column:header_img_url" json:"header_img_url"`
	HeaderVideoURL string `gorm:"column:header_video_url" json:"header_video_url"` // 建議修正為 HeaderVideoURL
	Conclusion     string `gorm:"column:conclusion" json:"conclusion"`
	Contents       string `gorm:"column:contents" json:"contents"`

	BodyImgURL  pq.StringArray `gorm:"column:body_img_url;type:text[]" json:"body_img_url"`
	BodyImgText pq.StringArray `gorm:"column:body_img_text;type:text[]" json:"body_img_text"`
	TitleTags   pq.StringArray `gorm:"column:title_tags;type:text[]" json:"title_tags"`
	AuthorNames pq.StringArray `gorm:"column:author_names;type:text[]" json:"author_names"`
	Editor      pq.StringArray `gorm:"column:editor;type:text[]" json:"editor"`
	Staff       pq.StringArray `gorm:"column:staff;type:text[]" json:"staff"`
	ContentTags pq.StringArray `gorm:"column:content_tags;type:text[]" json:"content_tags"`

	// 可為 NULL 的時間欄位 (PostgreSQL timestamptz -> Go *time.Time)
	FirstReleaseTime *time.Time `gorm:"column:first_release_time" json:"first_release_time,omitempty"`
	ReleaseTime      *time.Time `gorm:"column:release_time" json:"release_time,omitempty"`
	UpdateTime       *time.Time `gorm:"column:update_time" json:"update_time,omitempty"`
	ScraperUpdate    *time.Time `gorm:"column:scraper_update" json:"scraper_update,omitempty"`
}
