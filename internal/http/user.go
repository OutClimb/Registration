package http

import (
	"encoding/json"
	"net/http"

	"github.com/OutClimb/Registration/internal/app"
	"github.com/gin-gonic/gin"
)

type tokenPublic struct {
	Reset bool   `json:"reset"`
	Token string `json:"token"`
}

type userPublic struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (u *userPublic) Publicize(user *app.UserInternal) {
	u.Username = user.Username
	u.Role = user.Role
	u.Name = user.Name
	u.Email = user.Email
}

func (h *httpLayer) createToken(c *gin.Context) {
	// Get the authentication data
	bodyAsByteArray, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve request body"})
		return
	}

	// Parse the authentication data
	jsonMap := make(map[string]string)
	err = json.Unmarshal(bodyAsByteArray, &jsonMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse request body"})
		return
	}

	// Validate the authentication data
	if len(jsonMap["username"]) == 0 || len(jsonMap["password"]) == 0 || len(jsonMap["password"]) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password. Username and Password must be greater than 0 characters in length and password must be less than 72 characters in length"})
		return
	}

	// Authenticate the user
	if user, err := h.app.AuthenticateUser(jsonMap["username"], jsonMap["password"]); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	} else if token, err := h.app.CreateToken(user, c.ClientIP()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create token"})
		return
	} else {
		c.JSON(http.StatusOK, tokenPublic{Reset: user.RequirePasswordReset, Token: token})
	}
}

func (h *httpLayer) getSelf(c *gin.Context) {
	userId := c.GetUint("user_id")

	user, err := h.app.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to find user"})
		return
	}

	userPublic := userPublic{}
	userPublic.Publicize(user)

	c.JSON(http.StatusOK, userPublic)
}
