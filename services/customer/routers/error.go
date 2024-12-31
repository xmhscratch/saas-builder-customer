package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error comment
func (route *RouteContext) Error(err error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("Error: %+v\n", err)
		// log.Panic(err)

		ctx.JSON(http.StatusOK, gin.H{
			"error": map[string]interface{}{
				"message": err.Error(),
			},
		})
	}
}
