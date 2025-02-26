package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *httpLayer) createToken(c *gin.Context) {
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

	// Validate the authentication data
	if len(jsonMap["username"]) == 0 || len(jsonMap["password"]) == 0 || len(jsonMap["password"]) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Authenticate the user
	if user, err := h.app.AuthenticateUser(jsonMap["username"], jsonMap["password"]); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	} else if token, err := h.app.CreateToken(user, c.ClientIP()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
