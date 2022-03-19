package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CristianCurteanu/url-shortener/pkg"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"github.com/urfave/cli/v2"
)

func main() {
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
	app.SetRouter(nil)

	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "create-mapping",
				Usage: "Creates a new url key mapping",
				Action: func(c *cli.Context) error {
					url := c.String("url")
					key := urls.CreateKey(url)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					defer cancel()

					mapping := urls.UrlMapping{
						Key: key,
						URL: url,
					}
					err := app.UrlsDAO.Add(ctx, mapping)
					if err != nil {
						return err
					}

					result, err := json.Marshal(mapping)
					if err != nil {
						return err
					}

					fmt.Println(string(result))

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Usage:    "URL that will be mapped to a key",
						Required: true,
					},
				},
			},
			{
				Name:  "get-mapping",
				Usage: "Fetches url key mapping",
				Action: func(c *cli.Context) error {
					key := c.String("key")
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					defer cancel()

					mapping, err := app.UrlsDAO.SearchById(ctx, key)
					if err != nil {
						return err
					}
					result, err := json.Marshal(mapping)
					if err != nil {
						return err
					}

					fmt.Println(string(result))

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Usage:    "Key of the stored URL",
						Required: true,
					},
				},
			},
			{
				Name:  "delete-mapping",
				Usage: "Fetches url key mapping",
				Action: func(c *cli.Context) error {
					key := c.String("key")
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					defer cancel()

					err := app.UrlsDAO.Delete(ctx, key)
					if err != nil {
						return err
					}

					fmt.Printf("Mapping `%s` is deleted\n", key)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Usage:    "Key of the stored URL",
						Required: true,
					},
				},
			},
		},
	}

	err = cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
