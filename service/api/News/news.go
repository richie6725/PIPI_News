package NewsApi

import (
	"News/service/controller/newsCtrl"
	newsDaoModel "News/service/dao/daoModels/news"
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

		Source string `json:"source" valid:"required"`
	}{}

	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body: " + err.Error()})
		return
	}

	boArgs := &boNews.CreateOuterLayerArgs{
		Source: model.Source,
		Query:  make([]*newsDaoModel.OuterLayer, 0, len(model.OuterLayer)),
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
