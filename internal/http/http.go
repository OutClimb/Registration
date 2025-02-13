package http

import (
	"net/http"
	"os"

	"github.com/OutClimb/Registration/internal/app"
	"github.com/gin-gonic/gin"
)

type httpLayer struct {
	engine *gin.Engine
	app    app.AppLayer
}

func New(appLayer app.AppLayer) *httpLayer {
	h := &httpLayer{
		engine: gin.New(),
		app:    appLayer,
	}

	h.setupFrontendRoutes()
	h.setupApiRoutes()

	return h
}

func (h *httpLayer) setupFrontendRoutes() {
	h.engine.GET("/form/:slug", h.getForm)

	// If no route is matched, redirect to the main page
	h.engine.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "https://outclimb.gay")
	})
}

func (h *httpLayer) setupApiRoutes() {
	api := h.engine.Group("/api/v1")
	{
		api.GET("/ping", h.GetPing)
		api.POST("/submission/:formSlug", h.createSubmission)
	}
}

func (h *httpLayer) Run() {
	if addr, addrExists := os.LookupEnv("LISTEN_ADDR"); !addrExists {
		h.engine.Run("0.0.0.0:8080")
	} else {
		h.engine.Run(addr)
	}
}
