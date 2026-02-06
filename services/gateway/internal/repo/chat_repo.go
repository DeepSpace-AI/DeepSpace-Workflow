package repo

import (
	"context"
	"errors"
	"time"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) ListConversations(ctx context.Context, orgID, projectID int64) ([]model.Conversation, error) {
	var conversations []model.Conversation
	if err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ?", orgID, projectID).
		Order("updated_at DESC").
		Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ChatRepo) ListStandaloneConversations(ctx context.Context, orgID int64) ([]model.Conversation, error) {
	var conversations []model.Conversation
	if err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id IS NULL", orgID).
		Order("updated_at DESC").
		Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ChatRepo) CreateConversation(ctx context.Context, orgID, projectID int64, title *string) (*model.Conversation, error) {
	conv := model.Conversation{
		OrgID:     orgID,
		ProjectID: &projectID,
		Title:     title,
	}
	if err := r.db.WithContext(ctx).Create(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *ChatRepo) CreateStandaloneConversation(ctx context.Context, orgID int64, title *string) (*model.Conversation, error) {
	conv := model.Conversation{
		OrgID:     orgID,
		ProjectID: nil,
		Title:     title,
	}
	if err := r.db.WithContext(ctx).Create(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *ChatRepo) GetConversation(ctx context.Context, orgID, conversationID int64) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.WithContext(ctx).
		Where("id = ? AND org_id = ?", conversationID, orgID).
		First(&conv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conv, nil
}

func (r *ChatRepo) ListMessages(ctx context.Context, conversationID int64) ([]model.Message, error) {
	var messages []model.Message
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatRepo) CreateMessage(ctx context.Context, conversationID int64, role, content string, modelName, traceID *string) (*model.Message, error) {
	msg := model.Message{
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		Model:          modelName,
		TraceID:        traceID,
	}
	if err := r.db.WithContext(ctx).Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *ChatRepo) TouchConversation(ctx context.Context, conversationID int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Conversation{}).
		Where("id = ?", conversationID).
		Updates(map[string]any{"updated_at": time.Now()}).Error
}

func (r *ChatRepo) UpdateConversationTitle(ctx context.Context, orgID, conversationID int64, title string) (*model.Conversation, error) {
	if err := r.db.WithContext(ctx).
		Model(&model.Conversation{}).
		Where("id = ? AND org_id = ?", conversationID, orgID).
		Updates(map[string]any{"title": title, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	return r.GetConversation(ctx, orgID, conversationID)
}

func (r *ChatRepo) DeleteConversation(ctx context.Context, orgID, conversationID int64) (bool, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND org_id = ?", conversationID, orgID).
		Delete(&model.Conversation{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
