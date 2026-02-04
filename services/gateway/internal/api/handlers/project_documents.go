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

func (h *ProjectDocumentHandler) List(c *gin.Context) {
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

func (h *ProjectDocumentHandler) Create(c *gin.Context) {
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

func (h *ProjectDocumentHandler) Get(c *gin.Context) {
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

func (h *ProjectDocumentHandler) Update(c *gin.Context) {
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

func (h *ProjectDocumentHandler) Delete(c *gin.Context) {
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
