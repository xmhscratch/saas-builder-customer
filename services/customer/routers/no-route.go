package routers

import (
	"github.com/gin-gonic/gin"
)

// NoRoute comment
func (route *RouteContext) NoRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	}
}
