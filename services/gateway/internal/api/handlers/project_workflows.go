package handlers

import (
	"net/http"
	"strconv"

	"deepspace/internal/service/projectworkflow"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ProjectWorkflowHandler struct {
	svc *projectworkflow.Service
}

func NewProjectWorkflowHandler(svc *projectworkflow.Service) *ProjectWorkflowHandler {
	return &ProjectWorkflowHandler{svc: svc}
}

func (h *ProjectWorkflowHandler) List(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	items, err := h.svc.ListByProject(c.Request.Context(), orgID, projectID)
	if err != nil {
		respondInternal(c, "failed to list workflows")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createWorkflowRequest struct {
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Steps       datatypes.JSON  `json:"steps"`
}

func (h *ProjectWorkflowHandler) Create(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req createWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), orgID, projectID, req.Name, req.Description, req.Steps)
	if err != nil {
		if err == projectworkflow.ErrInvalidName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}
		respondInternal(c, "failed to create workflow")
		return
	}

	c.JSON(http.StatusCreated, item)
}

type updateWorkflowRequest struct {
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Steps       *datatypes.JSON `json:"steps"`
}

func (h *ProjectWorkflowHandler) Update(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	workflowID, err := strconv.ParseInt(c.Param("workflowId"), 10, 64)
	if err != nil || workflowID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workflow id"})
		return
	}

	var req updateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Update(c.Request.Context(), orgID, projectID, workflowID, req.Name, req.Description, req.Steps)
	if err != nil {
		switch err {
		case projectworkflow.ErrInvalidName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		case projectworkflow.ErrNoUpdates:
			c.JSON(http.StatusBadRequest, gin.H{"error": "no updates"})
			return
		default:
			respondInternal(c, "failed to update workflow")
			return
		}
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ProjectWorkflowHandler) Delete(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	workflowID, err := strconv.ParseInt(c.Param("workflowId"), 10, 64)
	if err != nil || workflowID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workflow id"})
		return
	}

	found, err := h.svc.Delete(c.Request.Context(), orgID, projectID, workflowID)
	if err != nil {
		respondInternal(c, "failed to delete workflow")
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
