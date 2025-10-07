package service

import (
	NewsApi "News/service/api/News"
	"News/service/controller/aclCtrl"
	"News/service/controller/newsCtrl"
	"News/service/controller/noteCtrl"
	"News/service/internal/config"
	"News/service/internal/database"
	"context"
	"fmt"
	"go.uber.org/dig"
	"net/http"
)

func News() Service {
	once.Do(func() {
		srv = &NewsServer{}
	})

	return srv
}

type NewsServer struct{}

func (srv *NewsServer) Run() {

	container := dig.New()
	srv.provideConfig(container)
	srv.provideService(container)
	srv.provideController(container)
	srv.provideCore(container)

	srv.invokeApiRoutes(container)

	if err := container.Invoke(srv.run); err != nil {
		panic(err)
	}

}

func (srv *NewsServer) provideConfig(container *dig.Container) {

	if err := container.Provide(config.NewNews); err != nil {
		panic(err)
	}
}

func (srv *NewsServer) provideService(container *dig.Container) {
	if err := container.Provide(func() context.Context {
		return context.TODO()
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(database.NewNews); err != nil {
		panic(err)
	}

	if err := container.Provide(NewsApi.NewServer); err != nil {
		panic(err)
	}

	if err := container.Provide(NewsApi.NewRouterRoot); err != nil {
		panic(err)
	}

	if err := container.Provide(NewsApi.NewGinEngine); err != nil {
		panic(err)
	}

}

func (srv *NewsServer) provideController(container *dig.Container) {
	if err := container.Provide(aclCtrl.NewAcl); err != nil {
		panic(err)
	}
	if err := container.Provide(noteCtrl.NewNote); err != nil {
		panic(err)
	}

	if err := container.Provide(newsCtrl.NewNews); err != nil {
		panic(err)
	}
}

func (srv *NewsServer) invokeApiRoutes(container *dig.Container) {
	if err := container.Invoke(NewsApi.NewServer); err != nil {
		panic(err)
	}

	if err := container.Invoke(NewsApi.NewGinEngine); err != nil {
		panic(err)
	}
	if err := container.Invoke(NewsApi.NewAcl); err != nil {
		panic(err)
	}
	if err := container.Invoke(NewsApi.NewNote); err != nil {
		panic(err)
	}

	if err := container.Invoke(NewsApi.NewNews); err != nil {
		panic(err)
	}

}

func (srv *NewsServer) provideCore(container *dig.Container) {}

func (srv *NewsServer) run(server *http.Server) {
	fmt.Printf("News starts at %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
