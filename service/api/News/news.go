package NewsApi

import (
	"News/service/controller/newsCtrl"
	newsDaoModel "News/service/dao/daoModels/news"
	daoModel "News/service/internal/database"
	boNews "News/service/internal/model/bo/news"
	"News/service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
)

func NewNews(pack newsApiPack) {
	c := &newsApi{pack: pack}
	group := pack.Root.Group("news")
	{
		group.POST("CreateNewsOuterLayer", c.createNewsOuterLayer)
		group.GET("GetNewsOuterLayer", c.getNewsOuterLayer)
		group.GET("GetFilterOuterLayer", c.getFilterOuterLayer)
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

func (api *newsApi) createNewsOuterLayer(ctx *gin.Context) {

	var model = struct {
		OuterLayer []struct {
			Title       string           `json:"title" valid:"required"`
			DetailURL   string           `json:"detail_url" valid:"required"`
			Tags        []string         `json:"tags" valid:"required"`
			ImgURL      string           `json:"img_url" valid:"required"`
			NewsID      int              `json:"news_id" valid:"required"`
			ReleaseTime utils.CustomTime `json:"release_time" valid:"required"`
			Category    string           `json:"category" valid:"required"`
			CategoryURL string           `json:"category_url" valid:"required"`
		} `json:"outer_layer"`

		SourceNews string `json:"source_news" valid:"required"`
	}{}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}

	boArgs := &boNews.CreateOuterLayerArgs{
		SourceNews: model.SourceNews,
		Query:      make([]*newsDaoModel.OuterLayer, 0, len(model.OuterLayer)),
	}

	for _, item := range model.OuterLayer {
		boArgs.Query = append(boArgs.Query, &newsDaoModel.OuterLayer{
			BaseNews: newsDaoModel.BaseNews{
				Title:       item.Title,
				DetailURL:   item.DetailURL,
				Tags:        item.Tags,
				ImgURL:      item.ImgURL,
				NewsID:      item.NewsID,
				ReleaseTime: item.ReleaseTime.Time,
				Category:    item.Category,
				CategoryURL: item.CategoryURL,
			},
		})
	}

	err := api.pack.NewsCtrl.CreateOuterLayer(ctx, boArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (api *newsApi) getNewsOuterLayer(ctx *gin.Context) {

	var model = struct {
		OuterLayer []struct {
			SourceNews string `json:"source_news" valid:"required"`
			NewsID     int    `json:"news_id" valid:"required"`
			Title      string `json:"title" valid:"required"`
		} `json:"outer_layer"`
	}{}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}

	boArgs := &boNews.GetOuterLayerArgs{
		Query: make([]*newsDaoModel.OuterLayer, 0, len(model.OuterLayer)),
	}

	for _, item := range model.OuterLayer {
		boArgs.Query = append(boArgs.Query, &newsDaoModel.OuterLayer{
			BaseNews: newsDaoModel.BaseNews{
				SourceNews: item.SourceNews,
				NewsID:     item.NewsID,
				Title:      item.Title,
			},
		})
	}

	reply, err := api.pack.NewsCtrl.GetOuterLayer(ctx, boArgs)
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

func (api *newsApi) getFilterOuterLayer(ctx *gin.Context) {
	var model = struct {
		Filter struct {
			SourceNews string   `json:"source_news"`
			Category   string   `json:"category"`
			Tags       []string `json:"tags" valid:"required"`
		} `json:"filter"`

		TimeInterval struct {
			StartTime utils.CustomTime `json:"start_time"`
			EndTime   utils.CustomTime `json:"end_time"`
		} `json:"time_interval"`

		Pagination struct {
			PageSize   int `json:"page_size"`
			PageNumber int `json:"page_number"`
		} `json:"pagination"`
	}{}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}

	boArgs := &boNews.GetFilterOuterLayerArgs{
		Query: &newsDaoModel.OuterLayer{
			BaseNews: newsDaoModel.BaseNews{
				SourceNews: model.Filter.SourceNews,
				Category:   model.Filter.Category,
				Tags:       model.Filter.Tags,
			},
		},

		TimeInterval: daoModel.TimeInterval{
			StartTime: model.TimeInterval.StartTime.Time,
			EndTime:   model.TimeInterval.EndTime.Time,
		},

		Pagination: daoModel.Pagenation{
			PageSize:   model.Pagination.PageSize,
			PageNumber: model.Pagination.PageNumber,
		},
	}

	reply, err := api.pack.NewsCtrl.GetFilterOuterLayer(ctx, boArgs)
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
