package http

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(h *httpLayer) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if expiration, err := claims.GetExpirationTime(); err != nil || !expiration.After(time.Now()) {
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		} else if userId, exists := claims["user_id"]; !exists {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		} else if err := h.app.ValidateToken(uint(userId.(float64)), c.ClientIP()); err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		} else {
			c.Set("user_id", userId)
			c.Next()
		}
	}
}
