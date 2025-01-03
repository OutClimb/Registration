package http

import (
	"fmt"
	"log"
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

	h.SetupFrontendRoutes()
	h.SetupApiRoutes()

	return h
}

func (h *httpLayer) SetupFrontendRoutes() {
	// Check if the web directory exists and is a directory
	stats, err := os.Stat("./web")
	if !os.IsNotExist(err) && stats.IsDir() {
		// Read the contents of the web directory
		entries, err := os.ReadDir("./web")
		if err != nil {
			log.Fatal(err)
		}

		// For each entry in the web directory, add a route
		for _, e := range entries {
			if e.IsDir() {
				h.engine.Static("/"+e.Name(), "./web/"+e.Name())
			} else {
				h.engine.StaticFile("/"+e.Name(), "./web/"+e.Name())
			}
		}
	} else {
		fmt.Println("The web directory does not exists, so no frontend will be served. Make sure to build the frontend first.")
	}

	// If no route is matched, serve the index.html file
	h.engine.NoRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File("./web/index.html")

		c.Status(http.StatusNotFound)
	})
}

func (h *httpLayer) SetupApiRoutes() {
	api := h.engine.Group("/api/v1")
	{
		api.GET("/ping", h.GetPing)
		api.GET("/event/:slug", h.GetEvent)
		api.POST("/register", h.CreateRegistration)
	}
}

func (h *httpLayer) Run() {
	if addr, addrExists := os.LookupEnv("LISTEN_ADDR"); !addrExists {
		h.engine.Run("0.0.0.0:8080")
	} else {
		h.engine.Run(addr)
	}
}
