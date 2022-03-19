package pkg

import (
	"context"
	"net/http"

	"github.com/CristianCurteanu/url-shortener/pkg/cache"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RedisConfig struct {
	URL      string
	Password string
}

type Config struct {
	Host     string
	MongoURL string
	Database string
	Redis    RedisConfig
}

type App struct {
	UrlsDAO urls.URLMappingDAO
	Cache   cache.UrlCache
	router  http.Handler
	Config  *Config
}

func NewApp(conf *Config) *App {
	return &App{Config: conf}
}

func (app *App) Init() error {
	err := app.SetupUrlsDAO()
	if err != nil {
		return err
	}

	app.SetupCache()

	app.SetupRouter()

	return nil
}

func (app *App) SetupCache() {
	app.Cache = cache.NewRedisCache(app.Config.Redis.URL, app.Config.Redis.Password)

}

func (app *App) SetupUrlsDAO() error {
	connection, err := mongo.Connect(context.Background(), options.Client().ApplyURI(app.Config.MongoURL))
	if err != nil {
		return err
	}

	db := connection.Database(app.Config.Database)
	app.UrlsDAO = urls.NewUrlMappingDao(db.Collection("mappings"))

	return nil
}

func (app *App) Router() http.Handler {
	return app.router
}

func (app *App) SetRouter(h http.Handler) *App {
	app.router = h
	return app
}
