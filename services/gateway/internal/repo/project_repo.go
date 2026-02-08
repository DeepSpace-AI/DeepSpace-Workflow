package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) ListByOrg(ctx context.Context, orgID int64) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", orgID).
		Order("id DESC").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepo) Get(ctx context.Context, orgID, projectID int64) (*model.Project, error) {
	var project model.Project
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", projectID, orgID).
		First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) Create(ctx context.Context, orgID int64, name string, description *string, projectType string) (*model.Project, error) {
	project := model.Project{
		UserID:      orgID,
		Name:        name,
		Type:        projectType,
		Description: description,
	}
	if err := r.db.WithContext(ctx).Create(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) Update(ctx context.Context, orgID, projectID int64, name *string, description *string) (*model.Project, error) {
	project, err := r.Get(ctx, orgID, projectID)
	if err != nil || project == nil {
		return project, err
	}

	updates := map[string]any{}
	if name != nil {
		updates["name"] = *name
	}
	if description != nil {
		updates["description"] = *description
	}

	if err := r.db.WithContext(ctx).Model(project).Updates(updates).Error; err != nil {
		return nil, err
	}

	if name != nil {
		project.Name = *name
	}
	if description != nil {
		project.Description = description
	}

	return project, nil
}

func (r *ProjectRepo) Delete(ctx context.Context, orgID, projectID int64) (bool, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", projectID, orgID).
		Delete(&model.Project{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ProjectRepo) CountByOrg(ctx context.Context, orgID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Project{}).
		Where("user_id = ?", orgID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
