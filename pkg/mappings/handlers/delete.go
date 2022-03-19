package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CristianCurteanu/url-shortener/pkg/mappings"
	"github.com/gin-gonic/gin"
)

func DeleteMappingHandler(urlMappingDAO mappings.URLMappingDAO) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := urlMappingDAO.Delete(ctx, id)
		if err != nil {
			log.Println("failed to delete mapping, error:", err)
			ctx.JSONP(http.StatusBadGateway, ErrorResponse{
				Key:     "mapping_not_deleted",
				Message: fmt.Sprintf("unable to remove the URL mapping for key `%s`", id),
			})

			return
		}

		ctx.JSONP(http.StatusOK, gin.H{"deleted": "OK"})
	}
}
