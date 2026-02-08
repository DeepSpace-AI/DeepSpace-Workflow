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

// List godoc
// @Summary 项目列表
// @Description 获取当前用户的项目列表
// @Tags 项目
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	items, err := h.svc.List(c.Request.Context(), orgID)
	if err != nil {
		respondInternal(c, "failed to list projects")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// Get godoc
// @Summary 获取项目详情
// @Description 获取指定项目详情
// @Tags 项目
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "项目不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// Create godoc
// @Summary 创建项目
// @Description 创建新项目
// @Tags 项目
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body createProjectRequest true "项目数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// Update godoc
// @Summary 更新项目
// @Description 更新项目名称或描述
// @Tags 项目
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param data body updateProjectRequest true "项目更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "项目不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id} [patch]
func (h *ProjectHandler) Update(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// Delete godoc
// @Summary 删除项目
// @Description 删除指定项目
// @Tags 项目
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "项目不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// Stats godoc
// @Summary 项目统计
// @Description 获取项目与资料统计
// @Tags 项目
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/stats [get]
func (h *ProjectHandler) Stats(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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
