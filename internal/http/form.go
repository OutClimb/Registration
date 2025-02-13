package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *httpLayer) getForm(c *gin.Context) {
	if !h.app.FormExists(c.Param("slug")) {
		c.Redirect(http.StatusTemporaryRedirect, "https://outclimb.gay")
	}

	if error := h.app.WriteFormTemplate(c.Param("slug"), c.Writer); error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to write form template",
		})
	}
}
