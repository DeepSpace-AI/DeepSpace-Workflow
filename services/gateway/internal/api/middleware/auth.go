package middleware

import (
	"errors"
	"net/http"
	"strings"

	"deepspace/internal/service/auth"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(validator *auth.APIKeyValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := extractAPIKey(c)
		if apiKey == "" {
			traceID, _ := c.Get("trace_id")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message":  "missing api key",
					"type":     "unauthorized",
					"trace_id": traceID,
				},
			})
			return
		}

		if validator == nil {
			traceID, _ := c.Get("trace_id")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message":  "auth validator not configured",
					"type":     "internal_error",
					"trace_id": traceID,
				},
			})
			return
		}

		authCtx, err := validator.Validate(c.Request.Context(), apiKey)
		if err != nil {
			traceID, _ := c.Get("trace_id")
			if errors.Is(err, auth.ErrInvalidAPIKey) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": gin.H{
						"message":  "invalid api key",
						"type":     "unauthorized",
						"trace_id": traceID,
					},
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message":  "auth validation failed",
					"type":     "internal_error",
					"trace_id": traceID,
				},
			})
			return
		}

		c.Set("api_key", apiKey)
		c.Set("api_key_id", authCtx.APIKeyID)
		c.Set("org_id", authCtx.OrgID)
		c.Next()
	}
}

func extractAPIKey(c *gin.Context) string {
	if auth := c.GetHeader("Authorization"); auth != "" {
		// Accept "Bearer <key>"
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return strings.TrimSpace(parts[1])
		}
	}

	if key := c.GetHeader("X-API-Key"); key != "" {
		return strings.TrimSpace(key)
	}

	return ""
}
