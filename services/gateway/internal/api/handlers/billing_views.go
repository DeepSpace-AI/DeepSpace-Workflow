package handlers

import (
	"net/http"
	"strconv"
	"time"

	"deepspace/internal/service/billing"
	"deepspace/internal/service/usage"

	"github.com/gin-gonic/gin"
)

type BillingViewHandler struct {
	billingSvc *billing.Service
	usageSvc   *usage.Service
}

func NewBillingViewHandler(billingSvc *billing.Service, usageSvc *usage.Service) *BillingViewHandler {
	return &BillingViewHandler{
		billingSvc: billingSvc,
		usageSvc:   usageSvc,
	}
}

// Wallet godoc
// @Summary 获取钱包
// @Description 获取当前用户钱包与近 24 小时用量
// @Tags 计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /billing/wallet [get]
func (h *BillingViewHandler) Wallet(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	wallet, err := h.billingSvc.GetWallet(c.Request.Context(), userID)
	if err != nil {
		respondInternal(c, "failed to load wallet")
		return
	}

	usage24h := float64(0)
	if h.usageSvc != nil {
		end := time.Now().UTC()
		start := end.Add(-24 * time.Hour)
		total, err := h.usageSvc.SumCost(c.Request.Context(), userID, &start, &end)
		if err != nil {
			respondInternal(c, "failed to load usage")
			return
		}
		usage24h = total
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet":    wallet,
		"usage_24h": usage24h,
	})
}

// Usage godoc
// @Summary 用量明细
// @Description 获取当前用户用量明细
// @Tags 计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param start query string false "开始时间（RFC3339）"
// @Param end query string false "结束时间（RFC3339）"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /billing/usage [get]
func (h *BillingViewHandler) Usage(c *gin.Context) {
	if h.usageSvc == nil {
		respondInternal(c, "usage service unavailable")
		return
	}
	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)

	var start *time.Time
	if value := c.Query("start"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start"})
			return
		}
		start = &parsed
	}

	var end *time.Time
	if value := c.Query("end"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end"})
			return
		}
		end = &parsed
	}

	if start != nil && end != nil && end.Before(*start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time range"})
		return
	}

	items, total, err := h.usageSvc.List(c.Request.Context(), usage.ListInput{
		UserID:   userID,
		Start:    start,
		End:      end,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		respondInternal(c, "failed to list usage")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

func parseIntQuery(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
