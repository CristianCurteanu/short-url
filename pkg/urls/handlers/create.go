package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/CristianCurteanu/url-shortener/pkg/cache"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"github.com/gin-gonic/gin"
)

type URLCreationRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type URLCreationResponse struct {
	Key string `json:"key"`
}

func CreateMappingHandler(urlMappingDAO urls.URLMappingDAO, cch cache.UrlCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req URLCreationRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSONP(http.StatusUnprocessableEntity, ErrorResponse{
				Key:     "unparsable_body",
				Message: err.Error(),
			})
			return
		}

		key := urls.CreateKey(req.URL)
		err = urlMappingDAO.Add(ctx, urls.UrlMapping{
			Key: key,
			URL: req.URL,
		})
		if err != nil {
			ctx.JSONP(http.StatusBadGateway, ErrorResponse{
				Key:     "failed_mapping_creation",
				Message: err.Error(),
			})
			return
		}

		go func(k, url string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*30)
			defer cancel()
			err := cch.Set(ctx, k, url)
			if err != nil {
				log.Printf("unable to store to cache, err: `%v`", err)
			}
		}(key, req.URL)

		ctx.JSONP(http.StatusCreated, URLCreationResponse{key})
	}
}
