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

// ListConversations godoc
// @Summary 项目对话列表
// @Description 获取指定项目的对话列表
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/conversations [get]
func (h *ChatSessionHandler) ListConversations(c *gin.Context) {
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

// ListStandaloneConversations godoc
// @Summary 独立对话列表
// @Description 获取独立对话列表
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /conversations [get]
func (h *ChatSessionHandler) ListStandaloneConversations(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	items, err := h.svc.ListStandaloneConversations(c.Request.Context(), orgID)
	if err != nil {
		respondInternal(c, "failed to list conversations")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// CreateConversation godoc
// @Summary 创建项目对话
// @Description 在指定项目中创建对话
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "项目ID"
// @Param data body createConversationRequest true "对话数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /projects/{id}/conversations [post]
func (h *ChatSessionHandler) CreateConversation(c *gin.Context) {
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

// CreateStandaloneConversation godoc
// @Summary 创建独立对话
// @Description 创建独立对话
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body createConversationRequest true "对话数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /conversations [post]
func (h *ChatSessionHandler) CreateStandaloneConversation(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	var req createConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item, err := h.svc.CreateStandaloneConversation(c.Request.Context(), orgID, req.Title)
	if err != nil {
		respondInternal(c, "failed to create conversation")
		return
	}

	c.JSON(http.StatusCreated, item)
}

// ListMessages godoc
// @Summary 对话消息列表
// @Description 获取对话消息列表
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param conversationId path int true "对话ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "对话不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /conversations/{conversationId}/messages [get]
func (h *ChatSessionHandler) ListMessages(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// CreateMessage godoc
// @Summary 创建对话消息
// @Description 在对话中追加消息
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param conversationId path int true "对话ID"
// @Param data body createMessageRequest true "消息数据"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "对话不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /conversations/{conversationId}/messages [post]
func (h *ChatSessionHandler) CreateMessage(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// UpdateConversation godoc
// @Summary 更新对话
// @Description 更新对话标题
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param conversationId path int true "对话ID"
// @Param data body updateConversationRequest true "对话更新数据"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "对话不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /conversations/{conversationId} [patch]
func (h *ChatSessionHandler) UpdateConversation(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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

// DeleteConversation godoc
// @Summary 删除对话
// @Description 删除指定对话
// @Tags 对话
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param conversationId path int true "对话ID"
// @Success 204 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "对话不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /conversations/{conversationId} [delete]
func (h *ChatSessionHandler) DeleteConversation(c *gin.Context) {
	orgID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
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
