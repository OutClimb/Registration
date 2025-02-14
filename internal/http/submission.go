package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *httpLayer) createSubmission(c *gin.Context) {
	form, error := h.app.GetForm(c.Param("slug"))
	if error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	bodyAsByteArray, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jsonMap := make(map[string]string)
	json.Unmarshal(bodyAsByteArray, &jsonMap)

	if err := h.app.ValidateSubmissionWithForm(jsonMap, form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if submission, err := h.app.CreateSubmission(c.Param("slug"), c.ClientIP(), c.GetHeader("User-Agent"), jsonMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, submission)
	}
}
