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
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header is required",
				"data":    nil,
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header format must be 'Bearer {token}'",
				"data":    nil,
			})
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Token has expired",
			"data":    nil,
		})
	case utils.ErrTokenInvalid, utils.ErrTokenMalformed:
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalid token",
			"data":    nil,
		})
	case utils.ErrTokenNotValidYet:
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Token not active yet",
			"data":    nil,
		})
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Authentication failed",
			"data":    nil,
		})
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
