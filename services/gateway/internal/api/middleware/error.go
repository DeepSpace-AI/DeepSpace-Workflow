package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		if c.Writer.Written() {
			return
		}

		traceID, _ := c.Get("trace_id")
		err := c.Errors.Last()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message":  err.Error(),
				"type":     "internal_error",
				"trace_id": traceID,
			},
		})
	}
}
