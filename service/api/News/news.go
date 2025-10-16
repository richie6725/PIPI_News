package NewsApi

import (
	"News/service/controller/newsCtrl"
	newsDaoModel "News/service/dao/daoModels/news"
	daoModel "News/service/internal/database"
	boNews "News/service/internal/model/bo/news"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
	"strings"
	"time"
)

func NewNews(pack newsApiPack) {
	c := &newsApi{pack: pack}
	group := pack.Root.Group("news")
	{
		group.POST("CreateNews", c.createNews)
		group.GET("GetNews", c.getNews)
		group.GET("GetNewsFilter", c.getNewsFilter)
	}

}

type newsApiPack struct {
	dig.In
	NewsCtrl newsCtrl.NewsCtrl
	Root     *gin.RouterGroup
}

type newsApi struct {
	pack newsApiPack
}

func (api *newsApi) createNews(ctx *gin.Context) {
	var model = struct {
		Data []struct {
			ID             int      `json:"id" valid:"required"`
			Title          string   `json:"title" valid:"required"`
			DetailURL      string   `json:"detail_url" valid:"required"`
			OutPageImgURL  string   `json:"out_page_img_url"`
			Category       string   `json:"category"`
			CategoryURL    string   `json:"category_url"`
			HeaderImgURL   string   `json:"header_img_url"`
			HeaderVideoURL string   `json:"header_vedio_url"`
			BodyImgURL     []string `json:"body_img_url"`
			BodyImgText    []string `json:"body_img_text"`
			Conclusion     string   `json:"conclusion"`
			Contents       []string `json:"contents"`

			TitleTags   string   `json:"title_tags"`
			AuthorNames []string `json:"author_names"`
			Editor      []string `json:"editor"`
			Staff       []string `json:"staff"`
			ContentTags []string `json:"content_tags"`

			// 時間欄位 (使用自訂 utils.CustomTime 型別)
			FirstReleaseTime int64 `json:"first_release_time"`
			ReleaseTime      int64 `json:"release_time"`
			UpdateTime       int64 `json:"update_time"`
			ScraperUpdate    int64 `json:"scraper_update"`
		} `json:"data"`
		SourceNews string `json:"source_news" valid:"required"`
	}{}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}

	boArgs := &boNews.CreateNewsArgs{
		SourceNews: model.SourceNews,
		Query:      make([]*newsDaoModel.News, 0, len(model.Data)),
	}

	boArgs.SourceNews = model.SourceNews

	for _, item := range model.Data {
		contentsStr := ""
		for _, content := range item.Contents {
			contentsStr += content + "\n" // 若不要換行可改成 += content
		}
		firstReleaseTime := time.UnixMilli(item.FirstReleaseTime)
		releaseTime := time.UnixMilli(item.ReleaseTime)
		updateTime := time.UnixMilli(item.UpdateTime)
		scraperUpdate := time.UnixMilli(item.ScraperUpdate)
		tagsSlice := strings.Split(item.TitleTags, ",")

		boArgs.Query = append(boArgs.Query, &newsDaoModel.News{
			ID:             item.ID,
			Title:          item.Title,
			DetailURL:      item.DetailURL,
			OutPageImgURL:  item.OutPageImgURL,
			Category:       item.Category,
			CategoryURL:    item.CategoryURL,
			HeaderImgURL:   item.HeaderImgURL,
			HeaderVideoURL: item.HeaderVideoURL,
			BodyImgURL:     item.BodyImgURL,
			BodyImgText:    item.BodyImgText,
			Conclusion:     item.Conclusion,
			Contents:       contentsStr,

			TitleTags:   tagsSlice,
			AuthorNames: item.AuthorNames,
			Editor:      item.Editor,
			Staff:       item.Staff,
			ContentTags: item.ContentTags,

			FirstReleaseTime: &firstReleaseTime,
			ReleaseTime:      &releaseTime,
			UpdateTime:       &updateTime,
			ScraperUpdate:    &scraperUpdate,
		})
	}

	err := api.pack.NewsCtrl.CreateNews(ctx, boArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

}

func (api *newsApi) getNews(ctx *gin.Context) {
	var model struct {
		SourceNews []string `json:"source_news" valid:"required"`
		ID         []int    `json:"id" valid:"required"`
	}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}
	if len(model.SourceNews) != len(model.ID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: source_news and id length not equal"})
		return
	}

	boArgs := &boNews.GetNewsArgs{
		Query:      make([]*newsDaoModel.News, 0, len(model.SourceNews)),
		SourceNews: make([]string, 0, len(model.ID)),
	}

	for i := range len(model.SourceNews) {
		boArgs.Query = append(boArgs.Query, &newsDaoModel.News{
			ID: model.ID[i],
		})
		boArgs.SourceNews = append(boArgs.SourceNews, model.SourceNews[i])
	}

	reply, err := api.pack.NewsCtrl.GetNews(ctx, boArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reply.Data) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "no data found",
		})
	}
	ctx.JSON(http.StatusOK, reply.Data)

}

func (api *newsApi) getNewsFilter(ctx *gin.Context) {
	var model = struct {
		Filter struct {
			SourceNews string   `json:"source_news"`
			Category   string   `json:"category"`
			Tags       []string `json:"tags" valid:"required"`
		} `json:"filter"`

		TimeInterval struct {
			StartTime time.Time `json:"start_time"`
			EndTime   time.Time `json:"end_time"`
		} `json:"time_interval"`

		Pagination struct {
			PageSize int `json:"page_size"`
			Page     int `json:"page"`
		} `json:"pagination"`
	}{}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}

	boArgs := &boNews.GetNewsFilterArgs{
		Query: &newsDaoModel.News{
			Category:  model.Filter.Category,
			TitleTags: model.Filter.Tags,
		},
		SourceNews: model.Filter.SourceNews,

		TimeInterval: &daoModel.TimeInterval{
			StartTime: model.TimeInterval.StartTime,
			EndTime:   model.TimeInterval.EndTime,
		},

		Pagination: &daoModel.Pagenation{
			PageSize: model.Pagination.PageSize,
			Page:     model.Pagination.Page,
		},
	}

	if boArgs.Pagination.PageSize > 20 {
		boArgs.Pagination.PageSize = 20
	}

	reply, err := api.pack.NewsCtrl.GetNewsFilter(ctx, boArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reply.Data) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "no data found",
		})
	}
	ctx.JSON(http.StatusOK, reply)

}
