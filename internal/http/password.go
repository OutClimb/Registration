package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *httpLayer) updatePassword(c *gin.Context) {
	userId := c.GetUint("user_id")
	user, err := h.app.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !user.RequirePasswordReset {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Get the authentication data
	bodyAsByteArray, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Error: Retrieving request body (%s)\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Parse the authentication data
	jsonMap := make(map[string]string)
	err = json.Unmarshal(bodyAsByteArray, &jsonMap)
	if err != nil {
		fmt.Printf("Error: Parsing request body (%s)\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal Server Error"})
		return
	}

	// Validate password
	if err := h.app.ValidatePassword(user, jsonMap["password"]); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the password
	if err := h.app.UpdatePassword(user, jsonMap["password"]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}
