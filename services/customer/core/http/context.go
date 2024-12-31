package http

import (
	"localdomain/customer/core"

	"github.com/gin-gonic/gin"
)

// Server comment
type Server struct {
	Config *core.Config
	Router *gin.Engine
}
