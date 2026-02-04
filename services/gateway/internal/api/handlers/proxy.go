package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"deepspace/internal/integrations/newapi"
	"deepspace/internal/pipeline"
	"deepspace/internal/pipeline/steps"
	"deepspace/internal/service/billing"
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
}

func NewProxyHandler(billingSvc *billing.Service, usageSvc *usage.Service, newapiClient *newapi.Client) *ProxyHandler {
	return &ProxyHandler{billing: billingSvc, usage: usageSvc, newapi: newapiClient}
}

func (h *ProxyHandler) Handle(c *gin.Context) {
	if h.newapi == nil {
		respondInternal(c, "newapi client not configured")
		return
	}

	orgID, ok := getOrgID(c)
	if !ok {
		respondInternal(c, "org_id missing")
		return
	}

	amount, hasAmount, err := parseBillingAmount(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing amount"})
		return
	}

	modelName := peekModelFromBody(c)

	state := pipeline.NewState()
	state.OrgID = orgID
	state.CostAmount = amount
	state.Model = modelName
	state.RefID = resolveRefID(c, hasAmount)
	if state.RefID == "" && hasAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ref_id is required for billing"})
		return
	}
	if value, ok := c.Get("api_key_id"); ok {
		if v, ok := value.(int64); ok {
			state.APIKeyID = v
		}
	}
	if value, ok := c.Get("trace_id"); ok {
		if v, ok := value.(string); ok {
			state.TraceID = v
		}
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

	post := pipeline.New(
		steps.NewUsageCapture(h.billing, h.usage),
	)
	_ = post.Run(c.Request.Context(), state)
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
