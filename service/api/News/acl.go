package NewsApi

import (
	"News/service/controller/aclCtrl"
	aclDaoModel "News/service/dao/daoModels/acl"
	boAcl "News/service/internal/model/bo/acl"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
	"time"
)

func NewAcl(pack aclApiPack) {
	c := &aclApi{pack: pack}
	group := pack.Root.Group("acl")
	{
		group.GET("status", c.getAcl)
		group.POST("status", c.updateAcl)
		group.GET("login", c.login)
	}

}

type aclApiPack struct {
	dig.In
	AclCtrl aclCtrl.AclCtrl
	Root    *gin.RouterGroup
}

type aclApi struct {
	pack aclApiPack
}

// 用於mongoDB抓使用者資料
func (api *aclApi) getAcl(ctx *gin.Context) {
	form := struct {
		Username string `json:"username" valid:"required"`
		Password string `json:"password" valid:"required"`
	}{}

	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	user := &boAcl.GetArgs{
		User: aclDaoModel.User{
			Username: form.Username,
			Password: form.Password,
		}}

	result, err := api.pack.AclCtrl.Get(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (api *aclApi) login(ctx *gin.Context) {
	form := struct {
		Username string `json:"username" valid:"required"`
		Password string `json:"password" valid:"required"`
		Token    string `json:"token" valid:"-"`
	}{}

	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	user := &boAcl.GetArgs{
		User: aclDaoModel.User{
			Username: form.Username,
			Password: form.Password,
			Token:    form.Token,
		}}

	result, err := api.pack.AclCtrl.GetLogin(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (api *aclApi) updateAcl(ctx *gin.Context) {
	form := struct {
		BulkUsers []aclDaoModel.BulkUserArg `json:"bulkUsers" valid:"required"`
	}{}

	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	user := &boAcl.UpdateArgs{
		Query: &aclDaoModel.Query{
			BulkUserArgs: form.BulkUsers,
			CreatedAt:    time.Now(),
		},
	}

	err := api.pack.AclCtrl.Update(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
