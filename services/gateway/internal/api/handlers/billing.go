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
	Amount   float64           `json:"amount"`
	RefID    string            `json:"ref_id"`
	Metadata map[string]any    `json:"metadata"`
}

func (h *BillingHandler) Hold(c *gin.Context) {
	h.handle(c, h.svc.Hold)
}

func (h *BillingHandler) Capture(c *gin.Context) {
	h.handle(c, h.svc.Capture)
}

func (h *BillingHandler) Release(c *gin.Context) {
	h.handle(c, h.svc.Release)
}

func (h *BillingHandler) handle(c *gin.Context, op func(ctx context.Context, orgID int64, amount float64, refID string, metadata map[string]any) (*billing.HoldResult, error)) {
	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
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

	result, err := op(c.Request.Context(), orgID, req.Amount, refID, req.Metadata)
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
