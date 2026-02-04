package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type ProjectDocumentRepo struct {
	db *gorm.DB
}

func NewProjectDocumentRepo(db *gorm.DB) *ProjectDocumentRepo {
	return &ProjectDocumentRepo{db: db}
}

func (r *ProjectDocumentRepo) ListByProject(ctx context.Context, orgID, projectID int64) ([]model.ProjectDocument, error) {
	var docs []model.ProjectDocument
	if err := r.db.WithContext(ctx).
		Where("org_id = ? AND project_id = ?", orgID, projectID).
		Order("updated_at DESC").
		Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *ProjectDocumentRepo) Get(ctx context.Context, orgID, projectID, docID int64) (*model.ProjectDocument, error) {
	var doc model.ProjectDocument
	err := r.db.WithContext(ctx).
		Where("id = ? AND org_id = ? AND project_id = ?", docID, orgID, projectID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *ProjectDocumentRepo) Create(ctx context.Context, doc *model.ProjectDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *ProjectDocumentRepo) Update(ctx context.Context, orgID, projectID, docID int64, updates map[string]any) (*model.ProjectDocument, error) {
	if err := r.db.WithContext(ctx).
		Model(&model.ProjectDocument{}).
		Where("id = ? AND org_id = ? AND project_id = ?", docID, orgID, projectID).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.Get(ctx, orgID, projectID, docID)
}

func (r *ProjectDocumentRepo) Delete(ctx context.Context, orgID, projectID, docID int64) (bool, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND org_id = ? AND project_id = ?", docID, orgID, projectID).
		Delete(&model.ProjectDocument{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
