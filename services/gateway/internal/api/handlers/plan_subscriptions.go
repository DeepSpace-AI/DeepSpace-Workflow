package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	planservice "deepspace/internal/service/plan"

	"github.com/gin-gonic/gin"
)

type PlanSubscriptionHandler struct {
	svc *planservice.Service
}

func NewPlanSubscriptionHandler(svc *planservice.Service) *PlanSubscriptionHandler {
	return &PlanSubscriptionHandler{svc: svc}
}

type subscriptionCreateRequest struct {
	UserID  int64   `json:"user_id"`
	PlanID  int64   `json:"plan_id"`
	Status  string  `json:"status"`
	StartAt string  `json:"start_at"`
	EndAt   *string `json:"end_at"`
}

type subscriptionUpdateRequest struct {
	Status  *string `json:"status"`
	StartAt *string `json:"start_at"`
	EndAt   *string `json:"end_at"`
}

func (h *PlanSubscriptionHandler) Create(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "订阅服务未配置")
		return
	}

	var req subscriptionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	startAt, err := parseRFC3339(req.StartAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间不正确"})
		return
	}
	var endAt *time.Time
	if req.EndAt != nil && strings.TrimSpace(*req.EndAt) != "" {
		value, err := parseRFC3339(*req.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不正确"})
			return
		}
		endAt = &value
	}

	userID := req.UserID
	if userID == 0 {
		if value, ok := getUserID(c); ok {
			userID = value
		}
	}
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}

	item, err := h.svc.CreateSubscription(c.Request.Context(), planservice.SubscriptionCreateInput{
		UserID:  userID,
		PlanID:  req.PlanID,
		Status:  strings.TrimSpace(req.Status),
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		handlePlanSubscriptionError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *PlanSubscriptionHandler) Update(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "订阅服务未配置")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订阅ID不正确"})
		return
	}

	var req subscriptionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	var startAt *time.Time
	if req.StartAt != nil && strings.TrimSpace(*req.StartAt) != "" {
		value, err := parseRFC3339(*req.StartAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间不正确"})
			return
		}
		startAt = &value
	}
	var endAt *time.Time
	if req.EndAt != nil && strings.TrimSpace(*req.EndAt) != "" {
		value, err := parseRFC3339(*req.EndAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不正确"})
			return
		}
		endAt = &value
	}

	item, err := h.svc.UpdateSubscription(c.Request.Context(), id, planservice.SubscriptionUpdateInput{
		Status:  req.Status,
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		handlePlanSubscriptionError(c, err)
		return
	}

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订阅不存在"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *PlanSubscriptionHandler) GetOrgActive(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "订阅服务未配置")
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}

	item, err := h.svc.GetActiveSubscription(c.Request.Context(), userID, time.Now().UTC())
	if err != nil {
		respondInternal(c, "获取订阅失败")
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订阅不存在"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func handlePlanSubscriptionError(c *gin.Context, err error) {
	switch err {
	case planservice.ErrInvalidSubscriptionTime:
		c.JSON(http.StatusBadRequest, gin.H{"error": "订阅时间不正确"})
	case planservice.ErrActiveSubscriptionExists:
		c.JSON(http.StatusConflict, gin.H{"error": "已有生效订阅"})
	case planservice.ErrPlanNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "套餐不存在"})
	default:
		respondInternal(c, "订阅操作失败")
	}
}

func parseRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.TrimSpace(value))
}
