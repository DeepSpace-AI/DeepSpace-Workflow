package repo

import (
	"context"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type ProjectWorkflowRepo struct {
	db *gorm.DB
}

func NewProjectWorkflowRepo(db *gorm.DB) *ProjectWorkflowRepo {
	return &ProjectWorkflowRepo{db: db}
}

func (r *ProjectWorkflowRepo) ListByProject(ctx context.Context, orgID, projectID int64) ([]model.ProjectWorkflow, error) {
	var items []model.ProjectWorkflow
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ?", orgID, projectID).
		Order("updated_at DESC").
		Find(&items).Error
	return items, err
}

func (r *ProjectWorkflowRepo) Get(ctx context.Context, orgID, projectID, workflowID int64) (*model.ProjectWorkflow, error) {
	var item model.ProjectWorkflow
	err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ? AND id = ?", orgID, projectID, workflowID).
		First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *ProjectWorkflowRepo) Create(ctx context.Context, item *model.ProjectWorkflow) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ProjectWorkflowRepo) Update(ctx context.Context, orgID, projectID, workflowID int64, updates map[string]any) (*model.ProjectWorkflow, error) {
	tx := r.db.WithContext(ctx).
		Model(&model.ProjectWorkflow{}).
		Where("org_id = ? AND project_id = ? AND id = ?", orgID, projectID, workflowID).
		Updates(updates)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return r.Get(ctx, orgID, projectID, workflowID)
}

func (r *ProjectWorkflowRepo) Delete(ctx context.Context, orgID, projectID, workflowID int64) (bool, error) {
	tx := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ? AND id = ?", orgID, projectID, workflowID).
		Delete(&model.ProjectWorkflow{})
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, nil
}
