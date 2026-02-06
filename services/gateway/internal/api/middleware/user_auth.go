package middleware

import (
	"net/http"
	"strings"

	"deepspace/internal/service/auth"

	"github.com/gin-gonic/gin"
)

func UserAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if jwtManager == nil {
			traceID, _ := c.Get("trace_id")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message":  "jwt manager not configured",
					"type":     "internal_error",
					"trace_id": traceID,
				},
			})
			return
		}

		token, err := c.Cookie(jwtManager.CookieName)
		if err != nil || token == "" {
			// Fallback for non-browser clients (optional).
			if authHeader := c.GetHeader("Authorization"); authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					token = parts[1]
				}
			}
		}
		if token == "" {
			traceID, _ := c.Get("trace_id")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message":  "missing session",
					"type":     "unauthorized",
					"trace_id": traceID,
				},
			})
			return
		}

		claims, err := jwtManager.Verify(token)
		if err != nil {
			traceID, _ := c.Get("trace_id")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message":  "invalid session",
					"type":     "unauthorized",
					"trace_id": traceID,
				},
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("org_id", claims.OrgID)
		c.Next()
	}
}
