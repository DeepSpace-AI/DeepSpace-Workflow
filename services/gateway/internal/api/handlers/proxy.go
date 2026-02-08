package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"deepspace/internal/integrations/newapi"
	"deepspace/internal/pipeline"
	"deepspace/internal/pipeline/steps"
	"deepspace/internal/service/billing"
	modelservice "deepspace/internal/service/model"
	planservice "deepspace/internal/service/plan"
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
}

func NewProxyHandler(billingSvc *billing.Service, usageSvc *usage.Service, newapiClient *newapi.Client, modelSvc *modelservice.Service, planSvc *planservice.Service) *ProxyHandler {
	return &ProxyHandler{billing: billingSvc, usage: usageSvc, newapi: newapiClient, model: modelSvc, plan: planSvc}
}

// Handle godoc
// @Summary 代理 NewAPI
// @Description 转发 /v1 下的请求到 NewAPI
// @Tags 代理
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param path path string true "转发路径"
// @Success 200 {object} map[string]interface{} "代理成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
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
		h.newapi.Proxy(c)
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
	modelID := ""

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
	if h.model != nil && modelName != "" {
		items, err := h.model.ListAll(c.Request.Context())
		if err == nil {
			for _, item := range items {
				if item.Name == modelName {
					modelID = item.ID
					state.Meta["price_input"] = item.PriceInput
					state.Meta["price_output"] = item.PriceOutput
					state.Meta["currency"] = item.Currency
					break
				}
			}
		}
	}
	if h.plan != nil {
		price, ok, err := h.plan.GetActivePlanPrice(c.Request.Context(), userID, modelID, time.Now().UTC())
		if err == nil && ok && price != nil {
			state.Meta["plan_applied"] = true
			state.Meta["plan_id"] = price.PlanID
			state.Meta["plan_billing_mode"] = price.BillingMode
			state.Meta["plan_currency"] = price.Currency
			state.Meta["plan_price_input"] = price.PriceInput
			state.Meta["plan_price_output"] = price.PriceOutput
			state.Meta["plan_price_request"] = price.PriceRequest
		}
	}
	if value, ok := c.Get("trace_id"); ok {
		if v, ok := value.(string); ok {
			state.TraceID = v
		}
	}
	if state.RefID == "" && state.TraceID != "" {
		state.RefID = state.TraceID
	}

	pre := pipeline.New(
		steps.NewAuth(),
		steps.NewPolicy(),
		steps.NewBudgetHold(h.billing),
	)
	if err := pre.Run(c.Request.Context(), state); err != nil {
		respondBillingError(c, err)
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
		steps.NewUsageCapture(h.billing, h.usage),
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
