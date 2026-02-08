package project

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/repo"
)

type Service struct {
	repo *repo.ProjectRepo
}

func New(repo *repo.ProjectRepo) *Service {
	return &Service{repo: repo}
}

var (
	ErrInvalidProjectName = errors.New("invalid project name")
	ErrNoProjectUpdates   = errors.New("no project updates")
)

type ProjectItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func (s *Service) List(ctx context.Context, userID int64) ([]ProjectItem, error) {
	items, err := s.repo.ListByOrg(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]ProjectItem, 0, len(items))
	for _, item := range items {
		result = append(result, ProjectItem{
			ID:          item.ID,
			Name:        item.Name,
			Type:        item.Type,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

func (s *Service) Get(ctx context.Context, userID, projectID int64) (*ProjectItem, error) {
	item, err := s.repo.Get(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	return &ProjectItem{
		ID:          item.ID,
		Name:        item.Name,
		Type:        item.Type,
		Description: item.Description,
		CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) Create(ctx context.Context, userID int64, name string, description *string, projectType string) (*ProjectItem, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Untitled Project"
	}

	if description != nil {
		value := strings.TrimSpace(*description)
		description = &value
	}

	projectType = normalizeProjectType(projectType)

	item, err := s.repo.Create(ctx, userID, name, description, projectType)
	if err != nil {
		return nil, err
	}

	return &ProjectItem{
		ID:          item.ID,
		Name:        item.Name,
		Type:        item.Type,
		Description: item.Description,
		CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) Update(ctx context.Context, userID, projectID int64, name *string, description *string) (*ProjectItem, error) {
	if name != nil {
		value := strings.TrimSpace(*name)
		if value == "" {
			return nil, ErrInvalidProjectName
		}
		name = &value
	}

	if description != nil {
		value := strings.TrimSpace(*description)
		description = &value
	}

	if name == nil && description == nil {
		return nil, ErrNoProjectUpdates
	}

	item, err := s.repo.Update(ctx, userID, projectID, name, description)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	return &ProjectItem{
		ID:          item.ID,
		Name:        item.Name,
		Type:        item.Type,
		Description: item.Description,
		CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) Delete(ctx context.Context, userID, projectID int64) (bool, error) {
	return s.repo.Delete(ctx, userID, projectID)
}

func (s *Service) CountByOrg(ctx context.Context, userID int64) (int64, error) {
	return s.repo.CountByOrg(ctx, userID)
}

func normalizeProjectType(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "research", "agent":
		return value
	default:
		return "research"
	}
}
