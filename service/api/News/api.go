package NewsApi

import (
	"News/service/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
	"time"
)

type servicePack struct {
	dig.In
	ServiceAddress config.ServiceAddress
	Handler        *gin.Engine
}

func NewServer(pack servicePack) *http.Server {
	return &http.Server{
		Addr:    pack.ServiceAddress.News,
		Handler: pack.Handler,
	}
}

func NewRouterRoot(pack servicePack) *gin.RouterGroup {
	return pack.Handler.Group("News")
}

func NewGinEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.ContextWithFallback = true
	router.Use(gin.Recovery(), cors.New(cors.Config{
		AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "X-Forwarded-For", "X-Real-IP"},
		MaxAge:          12 * time.Hour,
		AllowAllOrigins: true,
	}))

	return router

}
