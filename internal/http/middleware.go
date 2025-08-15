package http

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RedirectClaims struct {
	jwt.RegisteredClaims
	Audience string     `json:"aud"`
	User     userPublic `json:"user"`
}

func AuthMiddleware(h *httpLayer, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &RedirectClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
		} else if claims, ok := token.Claims.(*RedirectClaims); !ok {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
		} else if claims.Audience != c.ClientIP() {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		} else if userId, err := strconv.Atoi(claims.Subject); err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		} else if err := h.app.ValidateUser(uint(userId)); err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		} else if !h.app.CheckRole(claims.User.Role, role) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		} else {
			c.Set("user_id", uint(userId))
			c.Next()
		}
	}
}
