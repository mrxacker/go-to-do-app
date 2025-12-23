package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mrxacker/go-to-do-app/internal/infrastructure/auth"
)

func JWTMiddleware(jwtSvc *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()

		// List of routes to skip
		skipPaths := map[string]bool{
			"/api/v1/users/login":    true,
			"/api/v1/users/register": true,
		}

		if skipPaths[path] {
			c.Next() // skip middleware
			return
		}

		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		claims, err := jwtSvc.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
