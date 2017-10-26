package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/assistant-codelabs/conf"
)

// AuthMiddleware checks if the Authorization header's token matches the given
// configured token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header missing"})
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if token != conf.C.Token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
