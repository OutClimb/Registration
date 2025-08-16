package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	Audience string     `json:"aud"`
	User     userPublic `json:"user"`
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
	user, err := h.app.AuthenticateUser(jsonMap["username"], jsonMap["password"])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userPublic := userPublic{}
	userPublic.Publicize(user)

	if token, err := CreateToken(user.ID, &userPublic, c.ClientIP()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create token"})
		return
	} else {
		c.String(http.StatusOK, token)
	}
}

func CreateToken(userId uint, user *userPublic, clientIp string) (string, error) {
	// Get the token lifespan
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	if err != nil {
		return "", err
	}

	// Create the Claims
	claims := JwtClaims{}
	claims.Issuer = "registration"
	claims.Subject = strconv.FormatUint(uint64(userId), 10)
	claims.Audience = clientIp
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(tokenLifespan)))
	claims.NotBefore = jwt.NewNumericDate(time.Now())
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.User = *user

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", errors.New("Failed to sign token")
	} else {
		return signedToken, nil
	}
}

func JwtMiddleware(h *httpLayer, role string, resetAllowed bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
		} else if claims, ok := token.Claims.(*JwtClaims); !ok {
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
		} else if claims.User.RequirePasswordReset && !resetAllowed {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		} else {
			c.Set("user_id", uint(userId))
			c.Next()
		}
	}
}
