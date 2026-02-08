package steps

import (
	"context"
	"encoding/json"

	"deepspace/internal/pipeline"
	"deepspace/internal/service/billing"
	"deepspace/internal/service/usage"
)

type UsageCapture struct {
	billing *billing.Service
	usage   *usage.Service
}

func NewUsageCapture(billingSvc *billing.Service, usageSvc *usage.Service) *UsageCapture {
	return &UsageCapture{billing: billingSvc, usage: usageSvc}
}

func (s *UsageCapture) Name() string {
	return "usage_capture"
}

func (s *UsageCapture) Run(ctx context.Context, state *pipeline.State) error {
	if !hasProvidedBillingAmount(state) && state.CostAmount <= 0 {
		if planCost := calculateCostFromPlan(state); planCost > 0 {
			state.CostAmount = planCost
		} else {
			state.CostAmount = calculateCostFromUsage(state)
		}
	}

	if state.RefID != "" && state.CostAmount > 0 && s.billing != nil {
		if state.StatusCode >= 200 && state.StatusCode < 400 {
			if !hasProvidedBillingAmount(state) {
				if _, err := s.billing.Hold(ctx, state.UserID, state.CostAmount, state.RefID, map[string]any{"source": "pipeline"}); err != nil {
					return nil
				}
			}
			_, _ = s.billing.Capture(ctx, state.UserID, state.CostAmount, state.RefID, map[string]any{"source": "pipeline"})
		} else if hasProvidedBillingAmount(state) {
			_, _ = s.billing.Release(ctx, state.UserID, state.CostAmount, state.RefID, map[string]any{"source": "pipeline"})
		}
	}

	if s.usage != nil {
		_ = s.usage.Record(ctx, usage.RecordInput{
			UserID:           state.UserID,
			ProjectID:        nil,
			Model:            state.Model,
			PromptTokens:     state.UsagePromptTokens,
			CompletionTokens: state.UsageCompletionTokens,
			TotalTokens:      state.UsageTotalTokens,
			Cost:             state.CostAmount,
			TraceID:          state.TraceID,
		})
	}

	return nil
}

func calculateCostFromUsage(state *pipeline.State) float64 {
	if state == nil {
		return 0
	}
	if state.UsagePromptTokens <= 0 && state.UsageCompletionTokens <= 0 {
		return 0
	}
	priceInput := getMetaFloat(state.Meta, "price_input")
	priceOutput := getMetaFloat(state.Meta, "price_output")
	if priceInput <= 0 && priceOutput <= 0 {
		return 0
	}
	promptCost := (float64(state.UsagePromptTokens) / 1_000_000) * priceInput
	completionCost := (float64(state.UsageCompletionTokens) / 1_000_000) * priceOutput
	return promptCost + completionCost
}

func calculateCostFromPlan(state *pipeline.State) float64 {
	if state == nil {
		return 0
	}
	if !hasPlanApplied(state) {
		return 0
	}
	mode := getMetaString(state.Meta, "plan_billing_mode")
	switch mode {
	case "request":
		price := getMetaFloat(state.Meta, "plan_price_request")
		if price <= 0 {
			return 0
		}
		return price
	case "token":
		if state.UsagePromptTokens <= 0 && state.UsageCompletionTokens <= 0 {
			return 0
		}
		priceInput := getMetaFloat(state.Meta, "plan_price_input")
		priceOutput := getMetaFloat(state.Meta, "plan_price_output")
		if priceInput <= 0 && priceOutput <= 0 {
			return 0
		}
		promptCost := (float64(state.UsagePromptTokens) / 1_000_000) * priceInput
		completionCost := (float64(state.UsageCompletionTokens) / 1_000_000) * priceOutput
		return promptCost + completionCost
	default:
		return 0
	}
}

func getMetaFloat(meta map[string]any, key string) float64 {
	if meta == nil {
		return 0
	}
	value, ok := meta[key]
	if !ok {
		return 0
	}
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		parsed, err := v.Float64()
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func getMetaString(meta map[string]any, key string) string {
	if meta == nil {
		return ""
	}
	value, ok := meta[key]
	if !ok {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func hasPlanApplied(state *pipeline.State) bool {
	if state == nil || state.Meta == nil {
		return false
	}
	value, ok := state.Meta["plan_applied"]
	if !ok {
		return false
	}
	flag, ok := value.(bool)
	return ok && flag
}

func hasProvidedBillingAmount(state *pipeline.State) bool {
	if state == nil || state.Meta == nil {
		return false
	}
	value, ok := state.Meta["billing_amount_provided"]
	if !ok {
		return false
	}
	provided, ok := value.(bool)
	return ok && provided
}
