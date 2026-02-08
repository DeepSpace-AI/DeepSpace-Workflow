package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func castToInt64(value any) (int64, bool) {
	switch id := value.(type) {
	case int64:
		return id, true
	case int:
		return int64(id), true
	case int32:
		return int64(id), true
	case uint:
		return int64(id), true
	case uint32:
		return int64(id), true
	case uint64:
		if id > ^uint64(0)>>1 {
			return 0, false
		}
		return int64(id), true
	case float64:
		return int64(id), true
	case float32:
		return int64(id), true
	default:
		return 0, false
	}
}

func getUserID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}

	return castToInt64(value)
}

func respondInternal(c *gin.Context, message string) {
	traceID, ok := c.Get("trace_id")
	if ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message":  message,
				"type":     "internal_error",
				"trace_id": traceID,
			},
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}
