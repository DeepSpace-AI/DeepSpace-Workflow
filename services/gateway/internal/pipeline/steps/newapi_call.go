package steps

import (
	"context"

	"deepspace/internal/pipeline"
)

type NewAPICall struct{}

func NewNewAPICall() *NewAPICall {
	return &NewAPICall{}
}

func (s *NewAPICall) Name() string {
	return "newapi_call"
}

func (s *NewAPICall) Run(ctx context.Context, state *pipeline.State) error {
	// Placeholder for direct call integration.
	return nil
}
