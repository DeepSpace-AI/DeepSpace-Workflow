package handlers

import (
	"net/http"
	"strconv"

	"deepspace/internal/service/projectdocument"

	"github.com/gin-gonic/gin"
)

type ProjectDocumentHandler struct {
	svc *projectdocument.Service
}

func NewProjectDocumentHandler(svc *projectdocument.Service) *ProjectDocumentHandler {
	return &ProjectDocumentHandler{svc: svc}
}

// List godoc
// @Summary 项目文档列表
// @Description 获取指定项目的文档列表
// @Tags 项目文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/documents [get]
func (h *ProjectDocumentHandler) List(c *gin.Context) {
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
		respondInternal(c, "failed to list documents")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createDocumentRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// Create godoc
// @Summary 创建项目文档
// @Description 在指定项目中创建文档
// @Tags 项目文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param data body createDocumentRequest true "文档数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/documents [post]
func (h *ProjectDocumentHandler) Create(c *gin.Context) {
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

	var req createDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), orgID, projectID, req.Title, req.Content, req.Tags)
	if err != nil {
		respondInternal(c, "failed to create document")
		return
	}

	c.JSON(http.StatusCreated, item)
}

// Get godoc
// @Summary 获取项目文档
// @Description 获取指定项目的文档详情
// @Tags 项目文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param docId path int true "文档ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "文档不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/documents/{docId} [get]
func (h *ProjectDocumentHandler) Get(c *gin.Context) {
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

	docID, err := strconv.ParseInt(c.Param("docId"), 10, 64)
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	item, err := h.svc.Get(c.Request.Context(), orgID, projectID, docID)
	if err != nil {
		respondInternal(c, "failed to get document")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

type updateDocumentRequest struct {
	Title   *string   `json:"title"`
	Content *string   `json:"content"`
	Tags    *[]string `json:"tags"`
}

// Update godoc
// @Summary 更新项目文档
// @Description 更新指定项目的文档内容
// @Tags 项目文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param docId path int true "文档ID"
// @Param data body updateDocumentRequest true "文档更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "文档不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/documents/{docId} [patch]
func (h *ProjectDocumentHandler) Update(c *gin.Context) {
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

	docID, err := strconv.ParseInt(c.Param("docId"), 10, 64)
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	var req updateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.Update(c.Request.Context(), orgID, projectID, docID, req.Title, req.Content, req.Tags)
	if err != nil {
		switch err {
		case projectdocument.ErrInvalidTitle:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
			return
		case projectdocument.ErrNoUpdates:
			c.JSON(http.StatusBadRequest, gin.H{"error": "no updates"})
			return
		default:
			respondInternal(c, "failed to update document")
			return
		}
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Delete godoc
// @Summary 删除项目文档
// @Description 删除指定项目的文档
// @Tags 项目文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param docId path int true "文档ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "文档不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/documents/{docId} [delete]
func (h *ProjectDocumentHandler) Delete(c *gin.Context) {
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

	docID, err := strconv.ParseInt(c.Param("docId"), 10, 64)
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	found, err := h.svc.Delete(c.Request.Context(), orgID, projectID, docID)
	if err != nil {
		respondInternal(c, "failed to delete document")
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
