package repo

import (
	"context"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type ProjectSkillRepo struct {
	db *gorm.DB
}

func NewProjectSkillRepo(db *gorm.DB) *ProjectSkillRepo {
	return &ProjectSkillRepo{db: db}
}

func (r *ProjectSkillRepo) ListByProject(ctx context.Context, orgID, projectID int64) ([]model.ProjectSkill, error) {
	var items []model.ProjectSkill
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ?", orgID, projectID).
		Order("updated_at DESC").
		Find(&items).Error
	return items, err
}

func (r *ProjectSkillRepo) Get(ctx context.Context, orgID, projectID, skillID int64) (*model.ProjectSkill, error) {
	var item model.ProjectSkill
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ? AND id = ?", orgID, projectID, skillID).
		First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *ProjectSkillRepo) Create(ctx context.Context, item *model.ProjectSkill) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ProjectSkillRepo) Update(ctx context.Context, orgID, projectID, skillID int64, updates map[string]any) (*model.ProjectSkill, error) {
	tx := r.db.WithContext(ctx).
		Model(&model.ProjectSkill{}).
		Where("org_id = ? AND project_id = ? AND id = ?", orgID, projectID, skillID).
		Updates(updates)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return r.Get(ctx, orgID, projectID, skillID)
}

func (r *ProjectSkillRepo) Delete(ctx context.Context, orgID, projectID, skillID int64) (bool, error) {
	tx := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ? AND id = ?", orgID, projectID, skillID).
		Delete(&model.ProjectSkill{})
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, nil
}
