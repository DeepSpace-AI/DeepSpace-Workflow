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
