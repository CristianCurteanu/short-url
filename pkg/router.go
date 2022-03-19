package pkg

import (
	"github.com/CristianCurteanu/url-shortener/pkg/urls/handlers"
	"github.com/gin-gonic/gin"
)

func (app *App) SetupRouter() *App {
	router := gin.Default()
	router.GET("/:key", handlers.RedirectHandler(app.UrlsDAO, app.Cache))

	api := router.Group("/api")
	{
		api.POST("/mappings", handlers.CreateMappingHandler(app.UrlsDAO, app.Cache))
		api.GET("/mappings/:id", handlers.GetMappingHandler(app.UrlsDAO, app.Cache))
		api.GET("/mappings/:id/redirects", handlers.GetMappingRedirectsCounterHandler(app.UrlsDAO))
		api.DELETE("/mappings/:id", handlers.DeleteMappingHandler(app.UrlsDAO))
	}

	app.router = router

	return app
}
