package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const projectIDHeader = "X-Project-Id"

func ProjectContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := int64(0)

		if headerValue := strings.TrimSpace(c.GetHeader(projectIDHeader)); headerValue != "" {
			if parsed, err := strconv.ParseInt(headerValue, 10, 64); err == nil && parsed > 0 {
				projectID = parsed
			}
		}

		if projectID == 0 {
			if queryValue := strings.TrimSpace(c.Query("project_id")); queryValue != "" {
				if parsed, err := strconv.ParseInt(queryValue, 10, 64); err == nil && parsed > 0 {
					projectID = parsed
				}
			}
		}

		if projectID == 0 && requestMayHaveJSON(c) {
			raw, err := io.ReadAll(c.Request.Body)
			if err == nil && len(raw) > 0 {
				var payload struct {
					ProjectID *int64 `json:"project_id"`
				}
				if json.Unmarshal(raw, &payload) == nil && payload.ProjectID != nil && *payload.ProjectID > 0 {
					projectID = *payload.ProjectID
				}
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(raw))
		}

		c.Set("project_id", projectID)
		c.Next()
	}
}

func requestMayHaveJSON(c *gin.Context) bool {
	if c == nil || c.Request == nil {
		return false
	}
	if c.Request.Body == nil {
		return false
	}
	if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
		return false
	}
	contentType := strings.ToLower(strings.TrimSpace(c.GetHeader("Content-Type")))
	return strings.HasPrefix(contentType, "application/json")
}
