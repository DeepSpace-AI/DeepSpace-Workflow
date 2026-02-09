package handlers

import (
	"encoding/json"
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

// List godoc
// @Summary 项目工作流列表
// @Description 获取指定项目的工作流列表
// @Tags 项目工作流
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/workflows [get]
func (h *ProjectWorkflowHandler) List(c *gin.Context) {
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
		respondInternal(c, "failed to list workflows")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createWorkflowRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Steps       []any   `json:"steps"`
}

// Create godoc
// @Summary 创建项目工作流
// @Description 在指定项目中创建工作流
// @Tags 项目工作流
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param data body createWorkflowRequest true "工作流数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/workflows [post]
func (h *ProjectWorkflowHandler) Create(c *gin.Context) {
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

	var req createWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	stepsBytes, err := json.Marshal(req.Steps)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid steps"})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), orgID, projectID, req.Name, req.Description, datatypes.JSON(stepsBytes))
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
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Steps       *[]any  `json:"steps"`
}

// Update godoc
// @Summary 更新项目工作流
// @Description 更新指定项目的工作流
// @Tags 项目工作流
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param workflowId path int true "工作流ID"
// @Param data body updateWorkflowRequest true "工作流更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "工作流不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/workflows/{workflowId} [patch]
func (h *ProjectWorkflowHandler) Update(c *gin.Context) {
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

	var steps *datatypes.JSON
	if req.Steps != nil {
		stepsBytes, err := json.Marshal(*req.Steps)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid steps"})
			return
		}
		converted := datatypes.JSON(stepsBytes)
		steps = &converted
	}

	item, err := h.svc.Update(c.Request.Context(), orgID, projectID, workflowID, req.Name, req.Description, steps)
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

// Delete godoc
// @Summary 删除项目工作流
// @Description 删除指定项目的工作流
// @Tags 项目工作流
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param workflowId path int true "工作流ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "工作流不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/workflows/{workflowId} [delete]
func (h *ProjectWorkflowHandler) Delete(c *gin.Context) {
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
