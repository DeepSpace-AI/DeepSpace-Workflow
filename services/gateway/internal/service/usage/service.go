package usage

import (
	"context"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

type Service struct {
	repo *repo.UsageRepo
}

func New(repo *repo.UsageRepo) *Service {
	return &Service{repo: repo}
}

type ListInput struct {
	UserID   int64
	Start    *time.Time
	End      *time.Time
	Page     int
	PageSize int
}

type RecordInput struct {
	UserID           int64
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
		UserID:           in.UserID,
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

func (s *Service) List(ctx context.Context, in ListInput) ([]model.UsageRecord, int64, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	records, err := s.repo.List(ctx, repo.UsageListFilter{
		UserID: in.UserID,
		Start:  in.Start,
		End:    in.End,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, in.UserID, in.Start, in.End)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (s *Service) SumCost(ctx context.Context, userID int64, start, end *time.Time) (float64, error) {
	return s.repo.SumCost(ctx, userID, start, end)
}
