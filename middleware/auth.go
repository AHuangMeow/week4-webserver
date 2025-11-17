package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"week4-webserver/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Authorization header format must be 'Bearer {token}'")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			handleTokenError(c, err)
			return
		}

		c.Set("username", claims.Username)
		c.Set("token", tokenString)

		c.Next()
	}
}

func handleTokenError(c *gin.Context, err error) {
	switch err {
	case utils.ErrTokenExpired:
		utils.Unauthorized(c, "Token has expired")
	case utils.ErrTokenInvalid, utils.ErrTokenMalformed:
		utils.Unauthorized(c, "Invalid token")
	case utils.ErrTokenNotValidYet:
		utils.Unauthorized(c, "Token not active yet")
	default:
		utils.Unauthorized(c, "Authentication failed")
	}
	c.Abort()
}

func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	strUsername, ok := username.(string)
	if !ok || strUsername == "" {
		return "", false
	}

	return strUsername, true
}

func GetTokenFromContext(c *gin.Context) (string, bool) {
	token, exists := c.Get("token")
	if !exists {
		return "", false
	}

	strToken, ok := token.(string)
	if !ok || strToken == "" {
		return "", false
	}

	return strToken, true
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ParseToken(tokenString)
		if err == nil {
			c.Set("username", claims.Username)
			c.Set("token", tokenString)
		}

		c.Next()
	}
}
