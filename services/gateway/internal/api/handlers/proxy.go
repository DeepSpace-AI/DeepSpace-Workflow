package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/integrations/newapi"
	"deepspace/internal/pipeline"
	"deepspace/internal/pipeline/steps"
	"deepspace/internal/service/billing"
	modelservice "deepspace/internal/service/model"
	planservice "deepspace/internal/service/plan"
	"deepspace/internal/service/risk"
	"deepspace/internal/service/usage"

	"github.com/gin-gonic/gin"
)

const (
	billingAmountHeader = "X-Billing-Amount"
	billingRefHeader    = "X-Billing-Ref-Id"
)

type ProxyHandler struct {
	billing *billing.Service
	usage   *usage.Service
	newapi  *newapi.Client
	model   *modelservice.Service
	plan    *planservice.Service
	risk    *risk.Service
}

func NewProxyHandler(billingSvc *billing.Service, usageSvc *usage.Service, riskSvc *risk.Service, newapiClient *newapi.Client, modelSvc *modelservice.Service, planSvc *planservice.Service) *ProxyHandler {
	return &ProxyHandler{billing: billingSvc, usage: usageSvc, risk: riskSvc, newapi: newapiClient, model: modelSvc, plan: planSvc}
}

// Handle godoc
// @Summary 代理 NewAPI
// @Description 转发 /v1 下的请求到 NewAPI（不支持 /v1/models，且会校验模型是否允许）
// @Tags 代理
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param path path string true "转发路径"
// @Success 200 {object} map[string]interface{} "代理成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 404 {object} map[string]interface{} "接口不存在"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 402 {object} map[string]interface{} "余额不足"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /v1/{path} [get]
// @Router /v1/{path} [post]
// @Router /v1/{path} [put]
// @Router /v1/{path} [patch]
// @Router /v1/{path} [delete]
func (h *ProxyHandler) Handle(c *gin.Context) {
	if h.newapi == nil {
		respondInternal(c, "newapi client not configured")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id 缺失")
		return
	}

	if isModelListRequest(c) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// If billing is enabled but no amount is provided, still guard zero-balance usage.
	if h.billing != nil {
		wallet, err := h.billing.GetWallet(c.Request.Context(), userID)
		if err != nil {
			respondInternal(c, "failed to load wallet")
			return
		}
		if wallet != nil && wallet.Balance <= 0 {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient balance"})
			return
		}
	}

	amount, hasAmount, err := parseBillingAmount(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing amount"})
		return
	}

	modelName := peekModelFromBody(c)

	state := pipeline.NewState()
	state.UserID = userID
	state.CostAmount = amount
	state.Model = modelName
	if hasAmount {
		state.Meta["billing_amount_provided"] = true
	}
	state.RefID = resolveRefID(c, hasAmount)
	if state.RefID == "" && hasAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ref_id is required for billing"})
		return
	}
	modelName = strings.TrimSpace(modelName)
	if modelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "model is required"})
		return
	}
	if h.model != nil {
		item, err := h.model.GetActiveByName(c.Request.Context(), modelName)
		if err != nil {
			respondInternal(c, "failed to load model")
			return
		}
		if item == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "model not allowed"})
			return
		}
		state.Meta["price_input"] = item.PriceInput
		state.Meta["price_output"] = item.PriceOutput
		state.Meta["currency"] = item.Currency
	}
	if value, ok := c.Get("trace_id"); ok {
		if v, ok := value.(string); ok {
			state.TraceID = v
		}
	}
	state.Meta["client_ip"] = c.ClientIP()
	if value, ok := c.Get("project_id"); ok {
		switch v := value.(type) {
		case int64:
			if v > 0 {
				state.ProjectID = &v
			}
		case int:
			if v > 0 {
				parsed := int64(v)
				state.ProjectID = &parsed
			}
		case float64:
			if v > 0 {
				parsed := int64(v)
				state.ProjectID = &parsed
			}
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil && parsed > 0 {
				state.ProjectID = &parsed
			}
		}
	}
	if state.RefID == "" && state.TraceID != "" {
		state.RefID = state.TraceID
	}

	pre := pipeline.New(
		steps.NewAuth(),
		steps.NewPolicy(h.risk, h.usage),
		steps.NewBudgetHold(h.billing),
	)
	if err := pre.Run(c.Request.Context(), state); err != nil {
		switch {
		case errors.Is(err, steps.ErrRiskIPDenied):
			c.JSON(http.StatusForbidden, gin.H{"error": "IP 已被限制"})
		case errors.Is(err, steps.ErrRiskRateLimited):
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁"})
		case errors.Is(err, steps.ErrRiskBudgetExceeded):
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "预算已超限"})
		default:
			respondBillingError(c, err)
		}
		return
	}

	h.newapi.Proxy(c)
	state.StatusCode = c.Writer.Status()
	if usage, ok := newapi.GetUsageFromContext(c); ok {
		state.UsagePromptTokens = usage.PromptTokens
		state.UsageCompletionTokens = usage.CompletionTokens
		state.UsageTotalTokens = usage.TotalTokens
	}

	post := pipeline.New(
		steps.NewUsageCapture(h.billing, h.usage, h.plan),
	)
	_ = post.Run(c.Request.Context(), state)
}

func isModelListRequest(c *gin.Context) bool {
	if c.Request == nil {
		return false
	}
	if c.Request.Method != http.MethodGet {
		return false
	}
	path := c.Param("path")
	if path == "" {
		path = strings.TrimPrefix(c.Request.URL.Path, "/v1")
	}
	return path == "/models" || strings.HasSuffix(path, "/models")
}

func parseBillingAmount(c *gin.Context) (float64, bool, error) {
	value := strings.TrimSpace(c.GetHeader(billingAmountHeader))
	if value == "" {
		return 0, false, nil
	}

	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, true, err
	}
	return amount, true, nil
}

func resolveRefID(c *gin.Context, hasAmount bool) string {
	if !hasAmount {
		return ""
	}
	refID := strings.TrimSpace(c.GetHeader(billingRefHeader))
	if refID != "" {
		return refID
	}
	if trace, ok := c.Get("trace_id"); ok {
		if v, ok := trace.(string); ok {
			return v
		}
	}
	return ""
}

func peekModelFromBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(raw))

	var payload struct {
		Model string `json:"model"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return payload.Model
}
