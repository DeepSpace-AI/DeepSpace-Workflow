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

// List godoc
// @Summary 项目技能列表
// @Description 获取指定项目的技能列表
// @Tags 项目技能
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/skills [get]
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

// Create godoc
// @Summary 创建项目技能
// @Description 在指定项目中创建技能
// @Tags 项目技能
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param data body createSkillRequest true "技能数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/skills [post]
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

// Update godoc
// @Summary 更新项目技能
// @Description 更新指定项目的技能信息
// @Tags 项目技能
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param skillId path int true "技能ID"
// @Param data body updateSkillRequest true "技能更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "技能不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/skills/{skillId} [patch]
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

// Delete godoc
// @Summary 删除项目技能
// @Description 删除指定项目的技能
// @Tags 项目技能
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param skillId path int true "技能ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "技能不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/skills/{skillId} [delete]
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
