package http

import (
	"github.com/gin-gonic/gin"
)

func (h *httpLayer) GetEvent(c *gin.Context) {
	exists := h.app.CheckEventExists(c.Param("slug"))

	c.JSON(200, gin.H{
		"exists": exists,
	})
}
