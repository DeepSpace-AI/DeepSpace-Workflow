package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/service/knowledge"

	"github.com/gin-gonic/gin"
)

type KnowledgeHandler struct {
	svc *knowledge.Service
}

func NewKnowledgeHandler(svc *knowledge.Service) *KnowledgeHandler {
	return &KnowledgeHandler{svc: svc}
}

// ListBases godoc
// @Summary 知识库列表
// @Description 获取知识库列表
// @Tags 知识库
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param scope query string false "范围"
// @Param project_id query int false "项目ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases [get]
func (h *KnowledgeHandler) ListBases(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	scope := strings.ToLower(strings.TrimSpace(c.Query("scope")))
	if scope == "" {
		scope = "all"
	}

	var projectID *int64
	if projectParam := strings.TrimSpace(c.Query("project_id")); projectParam != "" {
		value, err := strconv.ParseInt(projectParam, 10, 64)
		if err != nil || value <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
			return
		}
		projectID = &value
	}

	items, err := h.svc.ListBases(c.Request.Context(), orgID, scope, projectID)
	if err != nil {
		respondInternal(c, "failed to list knowledge bases")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createKnowledgeBaseRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Scope       string  `json:"scope"`
	ProjectID   *int64  `json:"project_id"`
}

// CreateBase godoc
// @Summary 创建知识库
// @Description 创建新的知识库
// @Tags 知识库
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body createKnowledgeBaseRequest true "知识库数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "项目不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases [post]
func (h *KnowledgeHandler) CreateBase(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	var req createKnowledgeBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.CreateBase(c.Request.Context(), orgID, req.Scope, req.Name, req.Description, req.ProjectID)
	if err != nil {
		if err == knowledge.ErrInvalidScope || err == knowledge.ErrInvalidName || err == knowledge.ErrProjectRequired {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == knowledge.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		respondInternal(c, "failed to create knowledge base")
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetBase godoc
// @Summary 获取知识库
// @Description 获取知识库详情
// @Tags 知识库
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "知识库不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id} [get]
func (h *KnowledgeHandler) GetBase(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.svc.GetBase(c.Request.Context(), orgID, id)
	if err != nil {
		respondInternal(c, "failed to get knowledge base")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "knowledge base not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

type updateKnowledgeBaseRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// UpdateBase godoc
// @Summary 更新知识库
// @Description 更新知识库名称或描述
// @Tags 知识库
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Param data body updateKnowledgeBaseRequest true "知识库更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "知识库不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id} [patch]
func (h *KnowledgeHandler) UpdateBase(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateKnowledgeBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.UpdateBase(c.Request.Context(), orgID, id, req.Name, req.Description)
	if err != nil {
		if err == knowledge.ErrInvalidName || err == knowledge.ErrNoUpdates {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		respondInternal(c, "failed to update knowledge base")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "knowledge base not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteBase godoc
// @Summary 删除知识库
// @Description 删除指定知识库
// @Tags 知识库
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "知识库不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id} [delete]
func (h *KnowledgeHandler) DeleteBase(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	updated, err := h.svc.DeleteBase(c.Request.Context(), orgID, id)
	if err != nil {
		respondInternal(c, "failed to delete knowledge base")
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "knowledge base not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ListDocuments godoc
// @Summary 知识库文档列表
// @Description 获取指定知识库的文档列表
// @Tags 知识库文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id}/documents [get]
func (h *KnowledgeHandler) ListDocuments(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || kbID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base id"})
		return
	}

	items, err := h.svc.ListDocuments(c.Request.Context(), orgID, kbID)
	if err != nil {
		respondInternal(c, "failed to list documents")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// CreateDocument godoc
// @Summary 上传知识库文档
// @Description 上传文件到指定知识库
// @Tags 知识库文档
// @Accept multipart/form-data
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Param file formData file true "上传文件"
// @Success 201 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "知识库不存在"
// @Failure 413 {object} map[string]interface{} "文件过大"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id}/documents [post]
func (h *KnowledgeHandler) CreateDocument(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || kbID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base id"})
		return
	}

	if maxBytes := h.svc.MaxUploadBytes(); maxBytes > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		if err := c.Request.ParseMultipartForm(maxBytes); err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))
	var contentTypePtr *string
	if contentType != "" {
		contentTypePtr = &contentType
	}
	size := file.Size
	var sizePtr *int64
	if size > 0 {
		sizePtr = &size
	}

	src, err := file.Open()
	if err != nil {
		respondInternal(c, "failed to read file")
		return
	}
	defer src.Close()

	item, err := h.svc.CreateDocument(c.Request.Context(), orgID, kbID, file.Filename, contentTypePtr, sizePtr, src)
	if err != nil {
		switch err {
		case knowledge.ErrFileTooLarge:
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
			return
		case knowledge.ErrInvalidMimeType, knowledge.ErrInvalidFile:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == knowledge.ErrKnowledgeNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "knowledge base not found"})
			return
		}
		respondInternal(c, "failed to upload document")
		return
	}

	c.JSON(http.StatusCreated, item)
}

// DeleteDocument godoc
// @Summary 删除知识库文档
// @Description 删除指定知识库的文档
// @Tags 知识库文档
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Param docId path int true "文档ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "文档不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id}/documents/{docId} [delete]
func (h *KnowledgeHandler) DeleteDocument(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || kbID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base id"})
		return
	}

	docID, err := strconv.ParseInt(c.Param("docId"), 10, 64)
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	deleted, err := h.svc.DeleteDocument(c.Request.Context(), orgID, kbID, docID)
	if err != nil {
		respondInternal(c, "failed to delete document")
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// DownloadDocument godoc
// @Summary 下载知识库文档
// @Description 下载指定知识库的文档
// @Tags 知识库文档
// @Accept json
// @Produce application/octet-stream
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "知识库ID"
// @Param docId path int true "文档ID"
// @Param disposition query string false "下载方式"
// @Success 200 {file} file "文件流"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "文档不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /knowledge-bases/{id}/documents/{docId}/download [get]
func (h *KnowledgeHandler) DownloadDocument(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	kbID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || kbID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid knowledge base id"})
		return
	}

	docID, err := strconv.ParseInt(c.Param("docId"), 10, 64)
	if err != nil || docID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	doc, err := h.svc.GetDocument(c.Request.Context(), orgID, kbID, docID)
	if err != nil {
		respondInternal(c, "failed to get document")
		return
	}
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	if doc.ContentType != nil && strings.TrimSpace(*doc.ContentType) != "" {
		c.Header("Content-Type", *doc.ContentType)
	}

	disposition := strings.ToLower(strings.TrimSpace(c.Query("disposition")))
	if disposition == "inline" {
		c.Header("Content-Disposition", "inline; filename=\""+doc.FileName+"\"")
		c.File(doc.StoragePath)
		return
	}

	c.FileAttachment(doc.StoragePath, doc.FileName)
}
