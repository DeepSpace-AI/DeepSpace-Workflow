package handlers

import (
	"context"
	"net/http"
	"strings"

	"deepspace/internal/service/billing"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	svc *billing.Service
}

func NewBillingHandler(svc *billing.Service) *BillingHandler {
	return &BillingHandler{svc: svc}
}

type billingRequest struct {
	Amount   float64        `json:"amount"`
	RefID    string         `json:"ref_id"`
	Metadata map[string]any `json:"metadata"`
}

// Hold godoc
// @Summary 预扣余额
// @Description 预扣用户余额
// @Tags 计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body billingRequest true "预扣信息"
// @Success 200 {object} map[string]interface{} "预扣成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 402 {object} map[string]interface{} "余额不足"
// @Failure 409 {object} map[string]interface{} "引用冲突"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /billing/hold [post]
func (h *BillingHandler) Hold(c *gin.Context) {
	h.handle(c, h.svc.Hold)
}

// Capture godoc
// @Summary 扣减余额
// @Description 确认扣减预扣余额
// @Tags 计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body billingRequest true "扣减信息"
// @Success 200 {object} map[string]interface{} "扣减成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 402 {object} map[string]interface{} "余额不足"
// @Failure 409 {object} map[string]interface{} "引用冲突"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /billing/capture [post]
func (h *BillingHandler) Capture(c *gin.Context) {
	h.handle(c, h.svc.Capture)
}

// Release godoc
// @Summary 释放预扣
// @Description 释放已预扣余额
// @Tags 计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body billingRequest true "释放信息"
// @Success 200 {object} map[string]interface{} "释放成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 409 {object} map[string]interface{} "引用冲突"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /billing/release [post]
func (h *BillingHandler) Release(c *gin.Context) {
	h.handle(c, h.svc.Release)
}

func (h *BillingHandler) handle(c *gin.Context, op func(ctx context.Context, orgID int64, amount float64, refID string, metadata map[string]any) (*billing.HoldResult, error)) {
	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	var req billingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	refID := strings.TrimSpace(req.RefID)
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ref_id is required"})
		return
	}
	if req.Metadata == nil {
		req.Metadata = map[string]any{}
	}

	result, err := op(c.Request.Context(), userID, req.Amount, refID, req.Metadata)
	if err != nil {
		respondBillingError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func respondBillingError(c *gin.Context, err error) {
	switch err {
	case billing.ErrInvalidAmount:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	case billing.ErrInsufficientBalance:
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient balance"})
		return
	case billing.ErrInsufficientFrozen:
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient frozen balance"})
		return
	case billing.ErrRefConflict:
		c.JSON(http.StatusConflict, gin.H{"error": "ref_id conflict"})
		return
	default:
		respondInternal(c, "billing failed")
	}
}
