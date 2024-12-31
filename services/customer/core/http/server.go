package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	"localdomain/customer/core"
	"localdomain/customer/routers"
)

// NewServer comment
func NewServer(cfg *core.Config) (srv *Server, err error) {
	router := gin.Default()
	router.MaxMultipartMemory = 32 << 20 // 32 MiB
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.RemoveExtraSlash = true
	// router.SetTrustedProxies([]string{"172.17.0.1", "172.18.0.1", "127.0.0.1"})

	// route, err := routers.NewContext(cfg, solrInterface)
	route, err := routers.NewContext(cfg)
	if err != nil {
		return nil, err
	}
	route.Init(router)

	srv = &Server{
		Config: cfg,
		Router: router,
	}

	return srv, err
}

// Start comment
func (ctx *Server) Start() (err error) {
	if ctx.Router == nil {
		return errors.New("server is uninitialized")
	}
	portNumber := core.BuildString(":", ctx.Config.Port)
	return ctx.Router.Run(portNumber)
}
