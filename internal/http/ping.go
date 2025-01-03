package http

import (
	"github.com/gin-gonic/gin"
)

func (h *httpLayer) GetPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
