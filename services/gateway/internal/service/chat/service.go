package chat

import (
	"context"
	"strings"

	"deepspace/internal/repo"
)

type Service struct {
	repo *repo.ChatRepo
}

func New(repo *repo.ChatRepo) *Service {
	return &Service{repo: repo}
}

type ConversationItem struct {
	ID        int64   `json:"id"`
	ProjectID int64   `json:"project_id"`
	Title     *string `json:"title"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type MessageItem struct {
	ID             int64   `json:"id"`
	ConversationID int64   `json:"conversation_id"`
	Role           string  `json:"role"`
	Content        string  `json:"content"`
	Model          *string `json:"model"`
	TraceID        *string `json:"trace_id"`
	CreatedAt      string  `json:"created_at"`
}

func (s *Service) ListConversations(ctx context.Context, orgID, projectID int64) ([]ConversationItem, error) {
	items, err := s.repo.ListConversations(ctx, orgID, projectID)
	if err != nil {
		return nil, err
	}

	result := make([]ConversationItem, 0, len(items))
	for _, item := range items {
		result = append(result, ConversationItem{
			ID:        item.ID,
			ProjectID: projectID,
			Title:     item.Title,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

func (s *Service) CreateConversation(ctx context.Context, orgID, projectID int64, title *string) (*ConversationItem, error) {
	if title != nil {
		value := strings.TrimSpace(*title)
		title = &value
	}

	item, err := s.repo.CreateConversation(ctx, orgID, projectID, title)
	if err != nil {
		return nil, err
	}

	return &ConversationItem{
		ID:        item.ID,
		ProjectID: projectID,
		Title:     item.Title,
		CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) ListMessages(ctx context.Context, orgID, conversationID int64) ([]MessageItem, error) {
	conv, err := s.repo.GetConversation(ctx, orgID, conversationID)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		return nil, nil
	}

	items, err := s.repo.ListMessages(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	result := make([]MessageItem, 0, len(items))
	for _, item := range items {
		result = append(result, MessageItem{
			ID:             item.ID,
			ConversationID: item.ConversationID,
			Role:           item.Role,
			Content:        item.Content,
			Model:          item.Model,
			TraceID:        item.TraceID,
			CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

func (s *Service) CreateMessage(ctx context.Context, orgID, conversationID int64, role, content string, modelName, traceID *string) (*MessageItem, error) {
	conv, err := s.repo.GetConversation(ctx, orgID, conversationID)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		return nil, nil
	}

	role = strings.TrimSpace(role)
	if role == "" {
		role = "user"
	}
	content = strings.TrimSpace(content)

	item, err := s.repo.CreateMessage(ctx, conversationID, role, content, modelName, traceID)
	if err != nil {
		return nil, err
	}
	_ = s.repo.TouchConversation(ctx, conversationID)

	return &MessageItem{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		Role:           item.Role,
		Content:        item.Content,
		Model:          item.Model,
		TraceID:        item.TraceID,
		CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
