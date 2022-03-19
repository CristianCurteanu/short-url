package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CristianCurteanu/url-shortener/pkg/cache"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type URLMappingResponse struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

func GetMappingHandler(urlMappingDAO urls.URLMappingDAO, cch cache.UrlCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		url, cacheErr := cch.Get(ctx, id)
		if cacheErr == nil {
			ctx.JSONP(http.StatusOK, URLMappingResponse{id, url})
			return
		}

		mapping, err := urlMappingDAO.SearchById(ctx, id)
		if err != nil {
			log.Print(err)
			ctx.JSONP(http.StatusNotFound, ErrorResponse{
				Key:     "mapping_not_found",
				Message: fmt.Sprintf("URL mapping for key `%s` is not found", id),
			})
			return
		}

		if errors.Is(cacheErr, redis.Nil) {
			go func(k, url string) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*30)
				defer cancel()
				err := cch.Set(ctx, k, url)
				log.Printf("unable to store to cache, err: `%s`", err)
			}(mapping.Key, mapping.URL)
		}

		ctx.JSONP(http.StatusOK, URLMappingResponse{mapping.Key, mapping.URL})
	}
}

func GetMappingRedirectsCounterHandler(urlMappingDAO urls.URLMappingDAO) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		mapping, err := urlMappingDAO.SearchById(ctx, id)
		if err != nil {
			log.Print(err)
			ctx.JSONP(http.StatusNotFound, ErrorResponse{
				Key:     "mapping_not_found",
				Message: fmt.Sprintf("URL mapping for key `%s` is not found", id),
			})
			return
		}

		ctx.JSONP(http.StatusOK, gin.H{"key": id, "counter": mapping.Counter})
	}
}
