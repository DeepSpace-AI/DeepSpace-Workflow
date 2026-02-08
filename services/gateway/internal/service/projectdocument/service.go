package projectdocument

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

type Service struct {
	repo *repo.ProjectDocumentRepo
}

func New(repo *repo.ProjectDocumentRepo) *Service {
	return &Service{repo: repo}
}

var (
	ErrInvalidTitle = errors.New("invalid title")
	ErrNoUpdates    = errors.New("no updates")
)

type DocumentItem struct {
	ID        int64    `json:"id"`
	ProjectID int64    `json:"project_id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func (s *Service) ListByProject(ctx context.Context, userID, projectID int64) ([]DocumentItem, error) {
	items, err := s.repo.ListByProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	result := make([]DocumentItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapDocumentItem(&item))
	}
	return result, nil
}

func (s *Service) Get(ctx context.Context, userID, projectID, docID int64) (*DocumentItem, error) {
	item, err := s.repo.Get(ctx, userID, projectID, docID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	result := mapDocumentItem(item)
	return &result, nil
}

func (s *Service) Create(ctx context.Context, userID, projectID int64, title string, content string, tags []string) (*DocumentItem, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "未命名文档"
	}

	normalizedTags := normalizeTags(tags)
	tagJSON, _ := json.Marshal(normalizedTags)

	doc := &model.ProjectDocument{
		UserID:    userID,
		ProjectID: projectID,
		Title:     title,
		Content:   content,
		Tags:      tagJSON,
		Status:    "draft",
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, err
	}
	result := mapDocumentItem(doc)
	return &result, nil
}

func (s *Service) Update(ctx context.Context, userID, projectID, docID int64, title *string, content *string, tags *[]string) (*DocumentItem, error) {
	updates := map[string]any{}

	if title != nil {
		value := strings.TrimSpace(*title)
		if value == "" {
			return nil, ErrInvalidTitle
		}
		updates["title"] = value
	}
	if content != nil {
		updates["content"] = *content
	}
	if tags != nil {
		normalizedTags := normalizeTags(*tags)
		tagJSON, _ := json.Marshal(normalizedTags)
		updates["tags"] = tagJSON
	}

	if len(updates) == 0 {
		return nil, ErrNoUpdates
	}

	item, err := s.repo.Update(ctx, userID, projectID, docID, updates)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	result := mapDocumentItem(item)
	return &result, nil
}

func (s *Service) Delete(ctx context.Context, userID, projectID, docID int64) (bool, error) {
	return s.repo.Delete(ctx, userID, projectID, docID)
}

func mapDocumentItem(doc *model.ProjectDocument) DocumentItem {
	tags := []string{}
	if len(doc.Tags) > 0 {
		_ = json.Unmarshal(doc.Tags, &tags)
	}
	return DocumentItem{
		ID:        doc.ID,
		ProjectID: doc.ProjectID,
		Title:     doc.Title,
		Content:   doc.Content,
		Tags:      tags,
		Status:    doc.Status,
		CreatedAt: doc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: doc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func normalizeTags(tags []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		value := strings.TrimSpace(tag)
		if value == "" {
			continue
		}
		lower := strings.ToLower(value)
		if _, ok := seen[lower]; ok {
			continue
		}
		seen[lower] = struct{}{}
		result = append(result, value)
	}
	return result
}
