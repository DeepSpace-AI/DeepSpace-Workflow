package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/service/chat"

	"github.com/gin-gonic/gin"
)

type ChatSessionHandler struct {
	svc *chat.Service
}

func NewChatSessionHandler(svc *chat.Service) *ChatSessionHandler {
	return &ChatSessionHandler{svc: svc}
}

func (h *ChatSessionHandler) ListConversations(c *gin.Context) {
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

	items, err := h.svc.ListConversations(c.Request.Context(), orgID, projectID)
	if err != nil {
		respondInternal(c, "failed to list conversations")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createConversationRequest struct {
	Title *string `json:"title"`
}

func (h *ChatSessionHandler) CreateConversation(c *gin.Context) {
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

	var req createConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.CreateConversation(c.Request.Context(), orgID, projectID, req.Title)
	if err != nil {
		respondInternal(c, "failed to create conversation")
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ChatSessionHandler) ListMessages(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	conversationID, err := strconv.ParseInt(c.Param("conversationId"), 10, 64)
	if err != nil || conversationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	items, err := h.svc.ListMessages(c.Request.Context(), orgID, conversationID)
	if err != nil {
		respondInternal(c, "failed to list messages")
		return
	}
	if items == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

type createMessageRequest struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Model   *string `json:"model"`
	TraceID *string `json:"trace_id"`
}

func (h *ChatSessionHandler) CreateMessage(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	conversationID, err := strconv.ParseInt(c.Param("conversationId"), 10, 64)
	if err != nil || conversationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req createMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	item, err := h.svc.CreateMessage(c.Request.Context(), orgID, conversationID, req.Role, content, req.Model, req.TraceID)
	if err != nil {
		respondInternal(c, "failed to create message")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

type updateConversationRequest struct {
	Title string `json:"title"`
}

func (h *ChatSessionHandler) UpdateConversation(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	conversationID, err := strconv.ParseInt(c.Param("conversationId"), 10, 64)
	if err != nil || conversationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req updateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.UpdateConversation(c.Request.Context(), orgID, conversationID, req.Title)
	if err != nil {
		if err == chat.ErrInvalidTitle {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid title"})
			return
		}
		respondInternal(c, "failed to update conversation")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ChatSessionHandler) DeleteConversation(c *gin.Context) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	conversationID, err := strconv.ParseInt(c.Param("conversationId"), 10, 64)
	if err != nil || conversationID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	found, err := h.svc.DeleteConversation(c.Request.Context(), orgID, conversationID)
	if err != nil {
		respondInternal(c, "failed to delete conversation")
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
