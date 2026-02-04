package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const (
	traceHeader = "X-Trace-Id"
	requestIDHeader = "X-Request-Id"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(traceHeader)
		if traceID == "" {
			traceID = c.GetHeader(requestIDHeader)
		}
		if traceID == "" {
			traceID = newTraceID()
		}

		c.Set("trace_id", traceID)
		c.Writer.Header().Set(traceHeader, traceID)
		c.Request.Header.Set(traceHeader, traceID)

		c.Next()
	}
}

func newTraceID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "trace_unknown"
	}
	return hex.EncodeToString(b[:])
}
