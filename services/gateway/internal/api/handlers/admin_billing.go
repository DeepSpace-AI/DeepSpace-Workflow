package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"deepspace/internal/service/billing"
	"deepspace/internal/service/usage"

	"github.com/gin-gonic/gin"
)

type AdminBillingHandler struct {
	billingSvc *billing.Service
	usageSvc   *usage.Service
}

func NewAdminBillingHandler(billingSvc *billing.Service, usageSvc *usage.Service) *AdminBillingHandler {
	return &AdminBillingHandler{billingSvc: billingSvc, usageSvc: usageSvc}
}

type adminWalletUser struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type adminWalletItem struct {
	User          adminWalletUser `json:"user"`
	Balance       float64         `json:"balance"`
	FrozenBalance float64         `json:"frozen_balance"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type adminTopUpRequest struct {
	UserID   int64          `json:"user_id"`
	Amount   float64        `json:"amount"`
	Currency string         `json:"currency"`
	RefID    string         `json:"ref_id"`
	Metadata map[string]any `json:"metadata"`
}

// Wallets godoc
// @Summary 管理员：钱包列表
// @Description 获取钱包列表
// @Tags 管理-计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param user_id query int false "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/billing/wallets [get]
func (h *AdminBillingHandler) Wallets(c *gin.Context) {
	if h == nil || h.billingSvc == nil {
		respondInternal(c, "计费服务未配置")
		return
	}

	userID, err := parseOptionalInt64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}

	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)

	items, total, err := h.billingSvc.ListWallets(c.Request.Context(), billing.WalletListInput{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		respondInternal(c, "获取钱包失败")
		return
	}

	result := make([]adminWalletItem, 0, len(items))
	for _, item := range items {
		result = append(result, adminWalletItem{
			User: adminWalletUser{
				ID:        item.UserID,
				Email:     item.Email,
				Status:    item.Status,
				Role:      item.Role,
				CreatedAt: item.UserCreatedAt,
			},
			Balance:       item.Balance,
			FrozenBalance: item.FrozenBalance,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     result,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// Transactions godoc
// @Summary 管理员：交易流水
// @Description 获取交易流水
// @Tags 管理-计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param user_id query int false "用户ID"
// @Param type query string false "类型（hold/capture/release）"
// @Param start query string false "开始时间（RFC3339）"
// @Param end query string false "结束时间（RFC3339）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/billing/transactions [get]
func (h *AdminBillingHandler) Transactions(c *gin.Context) {
	if h == nil || h.billingSvc == nil {
		respondInternal(c, "计费服务未配置")
		return
	}

	userID, err := parseOptionalInt64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}

	start, end, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间范围不正确"})
		return
	}

	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)
	typeFilter := strings.TrimSpace(c.Query("type"))

	items, total, err := h.billingSvc.ListTransactions(c.Request.Context(), billing.TransactionListInput{
		UserID:   userID,
		Type:     typeFilter,
		Start:    start,
		End:      end,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		respondInternal(c, "获取流水失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// Usage godoc
// @Summary 管理员：用量记录
// @Description 获取用量记录
// @Tags 管理-计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param user_id query int false "用户ID"
// @Param start query string false "开始时间（RFC3339）"
// @Param end query string false "结束时间（RFC3339）"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/billing/usage [get]
func (h *AdminBillingHandler) Usage(c *gin.Context) {
	if h == nil || h.usageSvc == nil {
		respondInternal(c, "用量服务未配置")
		return
	}

	userID, err := parseOptionalInt64(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}

	start, end, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间范围不正确"})
		return
	}

	page := parseIntQueryAdmin(c, "page", 1)
	pageSize := parseIntQueryAdmin(c, "page_size", 20)

	items, total, err := h.usageSvc.ListAdmin(c.Request.Context(), usage.AdminListInput{
		UserID:   userID,
		Start:    start,
		End:      end,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		respondInternal(c, "获取用量失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// TopUp godoc
// @Summary 管理员：系统充值
// @Description 管理员给指定用户充值或冲正余额
// @Tags 管理-计费
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body adminTopUpRequest true "充值信息"
// @Success 200 {object} map[string]interface{} "充值成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 409 {object} map[string]interface{} "引用冲突"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/billing/topups [post]
func (h *AdminBillingHandler) TopUp(c *gin.Context) {
	if h == nil || h.billingSvc == nil {
		respondInternal(c, "计费服务未配置")
		return
	}

	var req adminTopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	if req.UserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不正确"})
		return
	}

	refID := strings.TrimSpace(req.RefID)
	if refID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ref_id 不能为空"})
		return
	}

	result, err := h.billingSvc.TopUp(c.Request.Context(), req.UserID, req.Amount, req.Currency, refID, req.Metadata)
	if err != nil {
		switch err {
		case billing.ErrInvalidAmount:
			c.JSON(http.StatusBadRequest, gin.H{"error": "金额不正确"})
		case billing.ErrRefConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "ref_id 冲突"})
		default:
			respondInternal(c, "充值失败")
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func parseOptionalInt64(value string) (*int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	if parsed <= 0 {
		return nil, errors.New("invalid user id")
	}
	return &parsed, nil
}

func parseTimeRange(c *gin.Context) (*time.Time, *time.Time, error) {
	var start *time.Time
	if value := strings.TrimSpace(c.Query("start")); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, nil, err
		}
		start = &parsed
	}
	var end *time.Time
	if value := strings.TrimSpace(c.Query("end")); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, nil, err
		}
		end = &parsed
	}
	if start != nil && end != nil && end.Before(*start) {
		return nil, nil, errors.New("invalid time range")
	}
	return start, end, nil
}

func parseIntQueryAdmin(c *gin.Context, key string, fallback int) int {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
