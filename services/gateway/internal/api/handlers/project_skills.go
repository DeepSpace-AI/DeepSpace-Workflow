package handlers

import (
	"net/http"
	"strconv"

	"deepspace/internal/service/projectskill"

	"github.com/gin-gonic/gin"
)

type ProjectSkillHandler struct {
	svc *projectskill.Service
}

func NewProjectSkillHandler(svc *projectskill.Service) *ProjectSkillHandler {
	return &ProjectSkillHandler{svc: svc}
}

func (h *ProjectSkillHandler) List(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	items, err := h.svc.ListByProject(c.Request.Context(), orgID, projectID)
	if err != nil {
		respondInternal(c, "failed to list skills")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createSkillRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Prompt      *string `json:"prompt"`
}

func (h *ProjectSkillHandler) Create(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req createSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), orgID, projectID, req.Name, req.Description, req.Prompt)
	if err != nil {
		if err == projectskill.ErrInvalidName {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		}
		respondInternal(c, "failed to create skill")
		return
	}

	c.JSON(http.StatusCreated, item)
}

type updateSkillRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Prompt      *string `json:"prompt"`
}

func (h *ProjectSkillHandler) Update(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	skillID, err := strconv.ParseInt(c.Param("skillId"), 10, 64)
	if err != nil || skillID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
		return
	}

	var req updateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Update(c.Request.Context(), orgID, projectID, skillID, req.Name, req.Description, req.Prompt)
	if err != nil {
		switch err {
		case projectskill.ErrInvalidName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
			return
		case projectskill.ErrNoUpdates:
			c.JSON(http.StatusBadRequest, gin.H{"error": "no updates"})
			return
		default:
			respondInternal(c, "failed to update skill")
			return
		}
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ProjectSkillHandler) Delete(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	skillID, err := strconv.ParseInt(c.Param("skillId"), 10, 64)
	if err != nil || skillID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
		return
	}

	found, err := h.svc.Delete(c.Request.Context(), orgID, projectID, skillID)
	if err != nil {
		respondInternal(c, "failed to delete skill")
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "skill not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
