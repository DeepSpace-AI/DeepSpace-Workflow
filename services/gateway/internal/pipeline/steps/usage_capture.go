package steps

import (
	"context"
	"encoding/json"
	"time"

	"deepspace/internal/pipeline"
	"deepspace/internal/service/billing"
	planservice "deepspace/internal/service/plan"
	"deepspace/internal/service/usage"
)

type UsageCapture struct {
	billing *billing.Service
	usage   *usage.Service
	plan    *planservice.Service
}

func NewUsageCapture(billingSvc *billing.Service, usageSvc *usage.Service, planSvc *planservice.Service) *UsageCapture {
	return &UsageCapture{billing: billingSvc, usage: usageSvc, plan: planSvc}
}

func (s *UsageCapture) Name() string {
	return "usage_capture"
}

func (s *UsageCapture) Run(ctx context.Context, state *pipeline.State) error {
	if !hasProvidedBillingAmount(state) && state.CostAmount <= 0 {
		state.CostAmount = calculateCostFromUsage(state)
	}

	if !hasProvidedBillingAmount(state) && state.StatusCode >= 200 && state.StatusCode < 400 && s.plan != nil {
		result, err := s.plan.ApplyQuota(ctx, state.UserID, time.Now().UTC(), int64(state.UsageTotalTokens), 1)
		if err == nil && result != nil && result.Applied {
			if result.Type == "request" {
				if result.Overage == 0 {
					state.CostAmount = 0
				}
			} else if result.Type == "token" {
				totalTokens := int64(state.UsageTotalTokens)
				if totalTokens > 0 {
					if result.Overage == 0 {
						state.CostAmount = 0
					} else if result.Overage < totalTokens {
						ratio := float64(result.Overage) / float64(totalTokens)
						state.CostAmount = state.CostAmount * ratio
					}
				}
			}
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
			ProjectID:        state.ProjectID,
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
