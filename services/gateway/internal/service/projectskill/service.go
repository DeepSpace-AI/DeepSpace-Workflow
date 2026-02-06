package projectskill

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

var (
	ErrInvalidName = errors.New("invalid name")
	ErrNoUpdates  = errors.New("no updates")
)

type Service struct {
	repo *repo.ProjectSkillRepo
}

func New(repo *repo.ProjectSkillRepo) *Service {
	return &Service{repo: repo}
}

type SkillItem struct {
	ID          int64   `json:"id"`
	OrgID       int64   `json:"org_id"`
	ProjectID   int64   `json:"project_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Prompt      *string `json:"prompt"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (s *Service) ListByProject(ctx context.Context, orgID, projectID int64) ([]SkillItem, error) {
	items, err := s.repo.ListByProject(ctx, orgID, projectID)
	if err != nil {
		return nil, err
	}
	result := make([]SkillItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapSkillItem(&item))
	}
	return result, nil
}

func (s *Service) Get(ctx context.Context, orgID, projectID, skillID int64) (*SkillItem, error) {
	item, err := s.repo.Get(ctx, orgID, projectID, skillID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	mapped := mapSkillItem(item)
	return &mapped, nil
}

func (s *Service) Create(ctx context.Context, orgID, projectID int64, name string, description *string, prompt *string) (*SkillItem, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}
	item := &model.ProjectSkill{
		OrgID:       orgID,
		ProjectID:   projectID,
		Name:        strings.TrimSpace(name),
		Description: description,
		Prompt:      prompt,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	mapped := mapSkillItem(item)
	return &mapped, nil
}

func (s *Service) Update(ctx context.Context, orgID, projectID, skillID int64, name *string, description *string, prompt *string) (*SkillItem, error) {
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
	if prompt != nil {
		updates["prompt"] = prompt
	}
	if len(updates) == 0 {
		return nil, ErrNoUpdates
	}
	item, err := s.repo.Update(ctx, orgID, projectID, skillID, updates)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	mapped := mapSkillItem(item)
	return &mapped, nil
}

func (s *Service) Delete(ctx context.Context, orgID, projectID, skillID int64) (bool, error) {
	return s.repo.Delete(ctx, orgID, projectID, skillID)
}

func mapSkillItem(item *model.ProjectSkill) SkillItem {
	return SkillItem{
		ID:          item.ID,
		OrgID:       item.OrgID,
		ProjectID:   item.ProjectID,
		Name:        item.Name,
		Description: item.Description,
		Prompt:      item.Prompt,
		CreatedAt:   item.CreatedAt.Format(timeFormat),
		UpdatedAt:   item.UpdatedAt.Format(timeFormat),
	}
}

const timeFormat = "2006-01-02 15:04:05"
