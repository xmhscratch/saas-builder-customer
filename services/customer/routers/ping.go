package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping comment
func (route *RouteContext) Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	}
}
