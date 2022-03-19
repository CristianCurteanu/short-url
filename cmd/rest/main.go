package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CristianCurteanu/url-shortener/pkg"
)

func main() {
	port := appPort()
	app := pkg.NewApp(&pkg.Config{
		Host:     appPort(),
		MongoURL: mongoURI(),
		Database: "url_mapping",
		Redis: pkg.RedisConfig{
			URL:      redisURI(),
			Password: redisPassword(),
		},
	})

	err := app.Init()
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    port,
		Handler: app.Router(),
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	log.Print("Listening on ", port)

	close := make(chan os.Signal)
	signal.Notify(close, os.Interrupt)
	<-close

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("server stopping...")
	defer cancel()

	log.Fatal(server.Shutdown(ctx))
	log.Fatal(app.Cache.Close())
}

func cacheDir() string {
	var err error
	dir := os.Getenv("CACHE_DIR")
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	return dir
}

func appPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT env var not defined")
	}
	return fmt.Sprintf(":%s", port)
}

func mongoURI() string {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		panic("MONGO_URI env var not defined")
	}
	return mongoURI
}

func redisURI() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		panic("REDIS_HOST env var is not defined")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		panic("REDIS_PORT env var is not defined")
	}

	return fmt.Sprintf("%s:%s", host, port)
}

func redisPassword() string {
	pass := os.Getenv("REDIS_PASSWORD")
	if pass == "" {
		panic("REDIS_PASSWORD env var is not defined")
	}

	return pass
}
