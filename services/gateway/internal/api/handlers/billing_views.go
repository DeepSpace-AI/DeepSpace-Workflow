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
