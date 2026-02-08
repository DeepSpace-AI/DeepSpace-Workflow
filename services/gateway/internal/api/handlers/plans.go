package handlers

import (
	"net/http"
	"strconv"
	"strings"

	planservice "deepspace/internal/service/plan"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	svc *planservice.Service
}

func NewPlanHandler(svc *planservice.Service) *PlanHandler {
	return &PlanHandler{svc: svc}
}

type planCreateRequest struct {
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	Currency     string  `json:"currency"`
	BillingMode  string  `json:"billing_mode"`
	PriceInput   float64 `json:"price_input"`
	PriceOutput  float64 `json:"price_output"`
	PriceRequest float64 `json:"price_request"`
}

type planUpdateRequest struct {
	Name         *string  `json:"name"`
	Status       *string  `json:"status"`
	Currency     *string  `json:"currency"`
	BillingMode  *string  `json:"billing_mode"`
	PriceInput   *float64 `json:"price_input"`
	PriceOutput  *float64 `json:"price_output"`
	PriceRequest *float64 `json:"price_request"`
}

func (h *PlanHandler) List(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "套餐服务未配置")
		return
	}

	items, err := h.svc.ListPlans(c.Request.Context())
	if err != nil {
		respondInternal(c, "获取套餐失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *PlanHandler) Create(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "套餐服务未配置")
		return
	}

	var req planCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	item, err := h.svc.CreatePlan(c.Request.Context(), planservice.PlanCreateInput{
		Name:         strings.TrimSpace(req.Name),
		Status:       strings.TrimSpace(req.Status),
		Currency:     strings.TrimSpace(req.Currency),
		BillingMode:  strings.TrimSpace(req.BillingMode),
		PriceInput:   req.PriceInput,
		PriceOutput:  req.PriceOutput,
		PriceRequest: req.PriceRequest,
	})
	if err != nil {
		handlePlanError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *PlanHandler) Update(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "套餐服务未配置")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "套餐ID不正确"})
		return
	}

	var req planUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	item, err := h.svc.UpdatePlan(c.Request.Context(), id, planservice.PlanUpdateInput{
		Name:         req.Name,
		Status:       req.Status,
		Currency:     req.Currency,
		BillingMode:  req.BillingMode,
		PriceInput:   req.PriceInput,
		PriceOutput:  req.PriceOutput,
		PriceRequest: req.PriceRequest,
	})
	if err != nil {
		handlePlanError(c, err)
		return
	}

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "套餐不存在"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func handlePlanError(c *gin.Context, err error) {
	switch err {
	case planservice.ErrInvalidPlanName:
		c.JSON(http.StatusBadRequest, gin.H{"error": "套餐名称不正确"})
	case planservice.ErrInvalidBillingMode:
		c.JSON(http.StatusBadRequest, gin.H{"error": "计费模式不正确"})
	case planservice.ErrInvalidPlanPrice:
		c.JSON(http.StatusBadRequest, gin.H{"error": "套餐价格不正确"})
	case planservice.ErrInvalidPlanCurrency:
		c.JSON(http.StatusBadRequest, gin.H{"error": "货币类型不正确"})
	case planservice.ErrPlanNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "套餐不存在"})
	default:
		respondInternal(c, "套餐操作失败")
	}
}
