package middleware

import (
	"net/http"

	"deepspace/internal/service/user"

	"github.com/gin-gonic/gin"
)

func RequireAdmin(userSvc *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Use the service to get user details including role
		// Note: We use GetMe because it retrieves the full user object with role
		userModel, _, _, err := userSvc.GetMe(c.Request.Context(), userID.(int64))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		if userModel.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		c.Next()
	}
}
