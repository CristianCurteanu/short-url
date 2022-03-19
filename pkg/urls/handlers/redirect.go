package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/CristianCurteanu/url-shortener/pkg/cache"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"github.com/gin-gonic/gin"
)

var ErrMappingNotFound = errors.New("mapping not found")

type ErrorResponse struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func RedirectHandler(urlMappingDAO urls.URLMappingDAO, cch cache.UrlCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")

		url, err := cch.Get(ctx, key)
		if err == nil {
			incrementRedirect(ctx, urlMappingDAO, key, url)
			return
		}

		mapping, err := urlMappingDAO.SearchById(ctx, key)
		if err != nil {
			log.Print(err)
			ctx.JSONP(http.StatusNotFound, ErrorResponse{
				Key:     "mapping_not_found",
				Message: fmt.Sprintf("URL mapping for key `%s` is not found", key),
			})
			return
		}

		err = urlMappingDAO.IncrementCounter(ctx, key)
		if err != nil {
			log.Print(err)
			ctx.JSONP(http.StatusBadGateway, ErrorResponse{
				Key:     "counter_not_incremented",
				Message: fmt.Sprintf("unable to increment counter for `%s` key", key),
			})
			return
		}

		incrementRedirect(ctx, urlMappingDAO, key, mapping.URL)
	}
}

func incrementRedirect(ctx *gin.Context, dao urls.URLMappingDAO, key, url string) {
	err := dao.IncrementCounter(ctx, key)
	if err != nil {
		log.Print(err)
		ctx.JSONP(http.StatusBadGateway, ErrorResponse{
			Key:     "counter_not_incremented",
			Message: fmt.Sprintf("unable to increment counter for `%s` key", key),
		})
		return
	}
	ctx.Redirect(http.StatusPermanentRedirect, url)
}
