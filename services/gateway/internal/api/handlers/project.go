package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/service/knowledge"
	"deepspace/internal/service/project"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	svc          *project.Service
	knowledgeSvc *knowledge.Service
}

func NewProjectHandler(svc *project.Service, knowledgeSvc *knowledge.Service) *ProjectHandler {
	return &ProjectHandler{svc: svc, knowledgeSvc: knowledgeSvc}
}

func (h *ProjectHandler) List(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	items, err := h.svc.List(c.Request.Context(), orgID)
	if err != nil {
		respondInternal(c, "failed to list projects")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ProjectHandler) Get(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.svc.Get(c.Request.Context(), orgID, projectID)
	if err != nil {
		respondInternal(c, "failed to get project")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

type createProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Type        string  `json:"type"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	name := strings.TrimSpace(req.Name)
	item, err := h.svc.Create(c.Request.Context(), orgID, name, req.Description, req.Type)
	if err != nil {
		respondInternal(c, "failed to create project")
		return
	}

	c.JSON(http.StatusCreated, item)
}

type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h *ProjectHandler) Update(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Update(c.Request.Context(), orgID, projectID, req.Name, req.Description)
	if err != nil {
		switch err {
		case project.ErrInvalidProjectName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		case project.ErrNoProjectUpdates:
			c.JSON(http.StatusBadRequest, gin.H{"error": "no updates"})
			return
		default:
			respondInternal(c, "failed to update project")
			return
		}
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	found, err := h.svc.Delete(c.Request.Context(), orgID, projectID)
	if err != nil {
		respondInternal(c, "failed to delete project")
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) Stats(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectCount, err := h.svc.CountByOrg(c.Request.Context(), orgID)
	if err != nil {
		respondInternal(c, "failed to count projects")
		return
	}

	literatureCount := int64(0)
	if h.knowledgeSvc != nil {
		count, err := h.knowledgeSvc.CountDocumentsByOrg(c.Request.Context(), orgID)
		if err != nil {
			respondInternal(c, "failed to count literature")
			return
		}
		literatureCount = count
	}

	c.JSON(http.StatusOK, gin.H{
		"projects":   projectCount,
		"literature": literatureCount,
		"workflows":  0,
		"agents":     0,
	})
}
