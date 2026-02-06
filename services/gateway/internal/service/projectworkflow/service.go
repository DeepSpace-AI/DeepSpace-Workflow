package projectworkflow

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"

	"gorm.io/datatypes"
)

var (
	ErrInvalidName = errors.New("invalid name")
	ErrNoUpdates  = errors.New("no updates")
)

type Service struct {
	repo *repo.ProjectWorkflowRepo
}

func New(repo *repo.ProjectWorkflowRepo) *Service {
	return &Service{repo: repo}
}

type WorkflowItem struct {
	ID          int64          `json:"id"`
	OrgID       int64          `json:"org_id"`
	ProjectID   int64          `json:"project_id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Steps       datatypes.JSON `json:"steps"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

func (s *Service) ListByProject(ctx context.Context, orgID, projectID int64) ([]WorkflowItem, error) {
	items, err := s.repo.ListByProject(ctx, orgID, projectID)
	if err != nil {
		return nil, err
	}
	result := make([]WorkflowItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapWorkflowItem(&item))
	}
	return result, nil
}

func (s *Service) Get(ctx context.Context, orgID, projectID, workflowID int64) (*WorkflowItem, error) {
	item, err := s.repo.Get(ctx, orgID, projectID, workflowID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	mapped := mapWorkflowItem(item)
	return &mapped, nil
}

func (s *Service) Create(ctx context.Context, orgID, projectID int64, name string, description *string, steps datatypes.JSON) (*WorkflowItem, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}
	item := &model.ProjectWorkflow{
		OrgID:       orgID,
		ProjectID:   projectID,
		Name:        strings.TrimSpace(name),
		Description: description,
		Steps:       steps,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	mapped := mapWorkflowItem(item)
	return &mapped, nil
}

func (s *Service) Update(ctx context.Context, orgID, projectID, workflowID int64, name *string, description *string, steps *datatypes.JSON) (*WorkflowItem, error) {
	updates := map[string]any{}
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return nil, ErrInvalidName
		}
		updates["name"] = strings.TrimSpace(*name)
	}
	if description != nil {
		updates["description"] = description
	}
	if steps != nil {
		updates["steps"] = *steps
	}
	if len(updates) == 0 {
		return nil, ErrNoUpdates
	}
	item, err := s.repo.Update(ctx, orgID, projectID, workflowID, updates)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	mapped := mapWorkflowItem(item)
	return &mapped, nil
}

func (s *Service) Delete(ctx context.Context, orgID, projectID, workflowID int64) (bool, error) {
	return s.repo.Delete(ctx, orgID, projectID, workflowID)
}

func mapWorkflowItem(item *model.ProjectWorkflow) WorkflowItem {
	return WorkflowItem{
		ID:          item.ID,
		OrgID:       item.OrgID,
		ProjectID:   item.ProjectID,
		Name:        item.Name,
		Description: item.Description,
		Steps:       item.Steps,
		CreatedAt:   item.CreatedAt.Format(timeFormat),
		UpdatedAt:   item.UpdatedAt.Format(timeFormat),
	}
}

const timeFormat = "2006-01-02 15:04:05"
