package middleware

import (
	"net/http"

	"github.com/A-Victory/e-commerce/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "missing authorization header"})
			return
		}

		claims, err := tokens.ValidateToken(clientToken)
		if err != "" {
			c.AbortWithStatusJSON(http, http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
