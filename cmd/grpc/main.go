package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/CristianCurteanu/url-shortener/pkg"
	gen "github.com/CristianCurteanu/url-shortener/pkg/grpc"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

type grpcMappingServer struct {
	gen.UnimplementedMappingsServiceServer
	app *pkg.App
}

func (s *grpcMappingServer) GetMapping(ctx context.Context, req *gen.GetMappingRequest) (*gen.GetMappingResponse, error) {
	url, cacheErr := s.app.Cache.Get(ctx, req.Key)
	if cacheErr == nil {
		return &gen.GetMappingResponse{Key: req.Key, Url: url}, nil
	}

	mapping, err := s.app.UrlsDAO.SearchById(ctx, req.Key)
	if err != nil {
		log.Println("failed search by ID:", err)
		return nil, err
	}

	if errors.Is(cacheErr, redis.Nil) {
		go func(k, url string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*30)
			defer cancel()
			err := s.app.Cache.Set(ctx, k, url)
			if err != nil {
				log.Printf("unable to store to cache, err: `%v`", err)
			}
		}(mapping.Key, mapping.URL)
	}

	return &gen.GetMappingResponse{Key: mapping.Key, Url: mapping.URL}, nil
}

func (s *grpcMappingServer) CreateMapping(ctx context.Context, req *gen.CreateMappingRequest) (*gen.CreateMappingResponse, error) {
	key := urls.CreateKey(req.GetUrl())
	err := s.app.UrlsDAO.Add(ctx, urls.UrlMapping{
		Key: key,
		URL: req.GetUrl(),
	})
	if err != nil {
		log.Println("failed to create mapping:", err)
		return nil, err
	}

	go func(k, url string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*30)
		defer cancel()
		err := s.app.Cache.Set(ctx, k, url)
		if err != nil {
			log.Printf("unable to store to cache, err: `%v`", err)
		}
	}(key, req.Url)

	return &gen.CreateMappingResponse{Key: key}, nil
}

func (s *grpcMappingServer) DeleteMapping(ctx context.Context, req *gen.DeleteMappingRequest) (*gen.DeleteMappingResponse, error) {
	err := s.app.UrlsDAO.Delete(ctx, req.Key)
	if err != nil {
		log.Println("failed to delete mapping, error:", err)
		return nil, fmt.Errorf("unable to remove the URL mapping for key `%s`", req.Key)
	}

	return &gen.DeleteMappingResponse{Deleted: "OK"}, nil
}

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

	conn, err := net.Listen("tcp", app.Config.Host)

	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}
	// Create grpc server
	grpcServer := grpc.NewServer()

	server := grpcMappingServer{
		app: app,
	}
	gen.RegisterMappingsServiceServer(grpcServer, &server)

	fmt.Println("Starting gRPC server at : ", app.Config.Host)
	if err := grpcServer.Serve(conn); err != nil {
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
