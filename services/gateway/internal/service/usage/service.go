package usage

import (
	"context"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

type Service struct {
	repo *repo.UsageRepo
}

func New(repo *repo.UsageRepo) *Service {
	return &Service{repo: repo}
}

type RecordInput struct {
	OrgID            int64
	ProjectID        *int64
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Cost             float64
	TraceID          string
}

func (s *Service) Record(ctx context.Context, in RecordInput) error {
	modelName := strings.TrimSpace(in.Model)
	if modelName == "" {
		modelName = "unknown"
	}

	rec := model.UsageRecord{
		OrgID:            in.OrgID,
		ProjectID:        in.ProjectID,
		Model:            modelName,
		PromptTokens:     in.PromptTokens,
		CompletionTokens: in.CompletionTokens,
		TotalTokens:      in.TotalTokens,
		Cost:             in.Cost,
		TraceID:          in.TraceID,
	}

	return s.repo.Create(ctx, rec)
}
