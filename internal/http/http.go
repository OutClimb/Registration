package http

import (
	"net/http"
	"os"
	"strings"

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

	if proxies, proxiesExist := os.LookupEnv("TRUSTED_PROXIES"); proxiesExist {
		h.engine.SetTrustedProxies(strings.Split(proxies, ","))
	}

	h.setupFrontendRoutes()
	h.setupApiRoutes()

	return h
}

func (h *httpLayer) setupFrontendRoutes() {
	// Load templates
	h.engine.LoadHTMLGlob("./web/*.html.tmpl")

	// Static Files
	h.engine.Static("/assets", "./web/assets")
	h.engine.StaticFile("/favicon.ico", "./web/favicon.ico")
	h.engine.StaticFile("/manifest.json", "./web/manifest.json")
	h.engine.StaticFile("/robots.txt", "./web/robots.txt")

	// Form Endpoint
	h.engine.GET("/form/:slug", h.getForm)

	// If no route is matched, redirect to the main page
	h.engine.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "https://outclimb.gay")
	})
}

func (h *httpLayer) setupApiRoutes() {
	api := h.engine.Group("/api/v1")
	{
		// Health Check
		api.GET("/ping", h.GetPing)

		// Form Endpoint
		api.GET("/form/:slug", h.getFormApi)

		// Form Submission Endpoint
		api.POST("/submission/:slug", h.createSubmission)

		// Login Route
		api.POST("/token", h.createToken)

		// User Authenticated Routes
		userApi := api.Group("/").Use(AuthMiddleware(h, "user", "api"))
		{
			userApi.GET("/self", h.getSelf)
		}

		userReset := api.Group("/").Use(AuthMiddleware(h, "user", "reset"))
		{
			userReset.PUT("/password", h.updatePassword)
		}

		// Viewer Authenticated Routes
		viewerApi := api.Group("/").Use(AuthMiddleware(h, "viewer", "api"))
		{
			viewerApi.GET("/form", h.getFormsApi)
			viewerApi.GET("/submission/:slug", h.getSubmissionsApi)
		}

		// Admin Authenticated Routes
		// admin := api.Group("/").Use(AuthMiddleware(h, "admin"))
		// {
		// 	admin.GET("/form", h.getForms)
		// 	admin.GET("/form/:slug", h.getForm)
		// 	admin.POST("/form", h.createForm)
		// 	admin.PUT("/form/:slug", h.updateForm)
		// 	admin.DELETE("/form/:slug", h.deleteForm)

		// 	admin.GET("/user", h.getUsers)
		// 	admin.GET("/user/:id", h.getUser)
		// 	admin.POST("/user", h.createUser)
		// 	admin.PUT("/user/:id", h.updateUser)
		// 	admin.DELETE("/user/:id", h.deleteUser)
		// }
	}
}

func (h *httpLayer) Run() {
	if addr, addrExists := os.LookupEnv("LISTEN_ADDR"); !addrExists {
		h.engine.Run("0.0.0.0:8080")
	} else {
		h.engine.Run(addr)
	}
}
