package steps

import (
	"context"

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
	if state.RefID != "" && state.CostAmount > 0 && s.billing != nil {
		if state.StatusCode >= 200 && state.StatusCode < 400 {
			_, _ = s.billing.Capture(ctx, state.OrgID, state.CostAmount, state.RefID, map[string]any{"source": "pipeline"})
		} else {
			_, _ = s.billing.Release(ctx, state.OrgID, state.CostAmount, state.RefID, map[string]any{"source": "pipeline"})
		}
	}

	if s.usage != nil {
		var apiKeyID *int64
		if state.APIKeyID > 0 {
			apiKeyID = &state.APIKeyID
		}
		_ = s.usage.Record(ctx, usage.RecordInput{
			OrgID:            state.OrgID,
			ProjectID:        nil,
			APIKeyID:         apiKeyID,
			Model:            state.Model,
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
			Cost:             state.CostAmount,
			TraceID:          state.TraceID,
		})
	}

	return nil
}
