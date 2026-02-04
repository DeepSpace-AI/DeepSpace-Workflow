package repo

import (
	"context"
	"errors"
	"time"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type APIKeyRepo struct {
	db *gorm.DB
}

func NewAPIKeyRepo(db *gorm.DB) *APIKeyRepo {
	return &APIKeyRepo{db: db}
}

func (r *APIKeyRepo) FindActiveByHash(ctx context.Context, keyHash string) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.WithContext(ctx).
		Where("key_hash = ? AND status = 'active'", keyHash).
		First(&key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &key, nil
}

func (r *APIKeyRepo) Create(ctx context.Context, orgID int64, name *string, keyHash, keyPrefix string, scopes []byte) (*model.APIKey, error) {
	key := model.APIKey{
		OrgID:     orgID,
		Name:      name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		Scopes:    scopes,
		Status:    "active",
	}

	if err := r.db.WithContext(ctx).Create(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *APIKeyRepo) ListByOrg(ctx context.Context, orgID int64) ([]model.APIKey, error) {
	var keys []model.APIKey
	if err := r.db.WithContext(ctx).
		Where("org_id = ?", orgID).
		Order("id DESC").
		Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *APIKeyRepo) Disable(ctx context.Context, orgID, keyID int64) (bool, error) {
	res := r.db.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("id = ? AND org_id = ? AND status = 'active'", keyID, orgID).
		Updates(map[string]any{"status": "disabled"})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *APIKeyRepo) Delete(ctx context.Context, orgID, keyID int64) (bool, error) {
	res := r.db.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("id = ? AND org_id = ? AND status != 'deleted'", keyID, orgID).
		Updates(map[string]any{"status": "deleted"})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *APIKeyRepo) UpdateScopes(ctx context.Context, orgID, keyID int64, scopes []byte) (bool, error) {
	res := r.db.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("id = ? AND org_id = ? AND status = 'active'", keyID, orgID).
		Updates(map[string]any{"scopes": scopes})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *APIKeyRepo) UpdateLastUsed(ctx context.Context, keyID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("id = ?", keyID).
		Updates(map[string]any{"last_used_at": &now}).Error
}
