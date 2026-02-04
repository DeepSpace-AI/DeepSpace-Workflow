package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/service/apikey"

	"github.com/gin-gonic/gin"
)

type APIKeyHandler struct {
	svc *apikey.Service
}

func NewAPIKeyHandler(svc *apikey.Service) *APIKeyHandler {
	return &APIKeyHandler{svc: svc}
}

type createKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

func (h *APIKeyHandler) Create(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	var req createKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	name := strings.TrimSpace(req.Name)
	scopes := normalizeScopes(req.Scopes)
	created, err := h.svc.Create(c.Request.Context(), orgID, name, scopes)
	if err != nil {
		respondInternal(c, "failed to create api key")
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *APIKeyHandler) List(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	items, err := h.svc.List(c.Request.Context(), orgID)
	if err != nil {
		respondInternal(c, "failed to list api keys")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *APIKeyHandler) Disable(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	updated, err := h.svc.Disable(c.Request.Context(), orgID, id)
	if err != nil {
		respondInternal(c, "failed to disable api key")
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "api key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "disabled"})
}

func (h *APIKeyHandler) UpdateScopes(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Scopes []string `json:"scopes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	scopes := normalizeScopes(req.Scopes)
	updated, err := h.svc.UpdateScopes(c.Request.Context(), orgID, id, scopes)
	if err != nil {
		respondInternal(c, "failed to update scopes")
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "api key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated", "scopes": scopes})
}

func (h *APIKeyHandler) Delete(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	updated, err := h.svc.Delete(c.Request.Context(), orgID, id)
	if err != nil {
		respondInternal(c, "failed to delete api key")
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "api key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func getOrgID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("org_id")
	if !ok {
		return 0, false
	}

	id, ok := value.(int64)
	if !ok {
		return 0, false
	}
	return id, true
}

func respondInternal(c *gin.Context, message string) {
	traceID, _ := c.Get("trace_id")
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"message":  message,
			"type":     "internal_error",
			"trace_id": traceID,
		},
	})
}

func normalizeScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return []string{}
	}
	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		value := strings.TrimSpace(scope)
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}
