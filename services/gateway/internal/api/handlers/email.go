package handlers

import (
	"errors"
	"net/http"
	"strings"

	"deepspace/internal/service/email"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	svc *email.Service
}

func NewEmailHandler(svc *email.Service) *EmailHandler {
	return &EmailHandler{svc: svc}
}

type sendEmailRequest struct {
	Type         string            `json:"type"`
	To           []string          `json:"to"`
	Subject      string            `json:"subject"`
	HTMLTemplate string            `json:"html_template"`
	TemplateData map[string]any    `json:"template_data"`
	HTML         string            `json:"html"`
	Text         string            `json:"text"`
	ReplyTo      string            `json:"reply_to"`
	Headers      map[string]string `json:"headers"`
}

type enqueueEmailRequest struct {
	Items []sendEmailRequest `json:"items"`
}

// Send godoc
// @Summary 发送邮件
// @Description 立即发送邮件
// @Tags 邮件
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body sendEmailRequest true "邮件内容"
// @Success 200 {object} map[string]interface{} "发送成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /email/send [post]
func (h *EmailHandler) Send(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "邮件服务未配置")
		return
	}

	var req sendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	input := email.EmailInput{
		Type:         strings.TrimSpace(req.Type),
		To:           normalizeRecipients(req.To),
		Subject:      strings.TrimSpace(req.Subject),
		HTMLTemplate: strings.TrimSpace(req.HTMLTemplate),
		TemplateData: req.TemplateData,
		HTML:         req.HTML,
		Text:         req.Text,
		ReplyTo:      strings.TrimSpace(req.ReplyTo),
		Headers:      req.Headers,
	}

	if err := h.svc.Send(c.Request.Context(), input); err != nil {
		switch {
		case errors.Is(err, email.ErrEmailDisabled):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件发送未启用"})
			return
		case errors.Is(err, email.ErrInvalidEmailType):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件类型不正确"})
			return
		case errors.Is(err, email.ErrInvalidRecipients):
			c.JSON(http.StatusBadRequest, gin.H{"error": "收件人为空"})
			return
		case errors.Is(err, email.ErrMissingSubject):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件主题为空"})
			return
		case errors.Is(err, email.ErrTemplateOverride):
			c.JSON(http.StatusBadRequest, gin.H{"error": "不允许覆盖模板内容"})
			return
		default:
			respondInternal(c, "邮件发送失败")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Enqueue godoc
// @Summary 邮件入队
// @Description 批量邮件入队
// @Tags 邮件
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body enqueueEmailRequest true "邮件列表"
// @Success 200 {object} map[string]interface{} "入队成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /email/enqueue [post]
func (h *EmailHandler) Enqueue(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "邮件服务未配置")
		return
	}

	var req enqueueEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮件列表为空"})
		return
	}

	inputs := make([]email.EmailInput, 0, len(req.Items))
	for _, item := range req.Items {
		input := email.EmailInput{
			Type:         strings.TrimSpace(item.Type),
			To:           normalizeRecipients(item.To),
			Subject:      strings.TrimSpace(item.Subject),
			HTMLTemplate: strings.TrimSpace(item.HTMLTemplate),
			TemplateData: item.TemplateData,
			HTML:         item.HTML,
			Text:         item.Text,
			ReplyTo:      strings.TrimSpace(item.ReplyTo),
			Headers:      item.Headers,
		}
		inputs = append(inputs, input)
	}

	if err := h.svc.EnqueueBatch(c.Request.Context(), inputs); err != nil {
		switch {
		case errors.Is(err, email.ErrQueueUnavailable):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件队列不可用"})
			return
		case errors.Is(err, email.ErrInvalidEmailType):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件类型不正确"})
			return
		case errors.Is(err, email.ErrInvalidRecipients):
			c.JSON(http.StatusBadRequest, gin.H{"error": "收件人为空"})
			return
		case errors.Is(err, email.ErrMissingSubject):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮件主题为空"})
			return
		case errors.Is(err, email.ErrTemplateOverride):
			c.JSON(http.StatusBadRequest, gin.H{"error": "不允许覆盖模板内容"})
			return
		default:
			respondInternal(c, "邮件入队失败")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"count":  len(inputs),
	})
}

func normalizeRecipients(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}
