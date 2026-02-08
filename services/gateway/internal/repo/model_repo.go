package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModelRepo struct {
	db *gorm.DB
}

type NameProvider struct {
	Name     string
	Provider string
}

func NewModelRepo(db *gorm.DB) *ModelRepo {
	return &ModelRepo{db: db}
}

func (r *ModelRepo) ListAll(ctx context.Context) ([]model.Model, error) {
	var models []model.Model
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *ModelRepo) ListActive(ctx context.Context) ([]model.Model, error) {
	var models []model.Model
	if err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *ModelRepo) ListActiveByProvider(ctx context.Context, provider string) ([]model.Model, error) {
	var models []model.Model
	if err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Where("provider = ?", provider).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *ModelRepo) ListByPairs(ctx context.Context, pairs []NameProvider) ([]model.Model, error) {
	if len(pairs) == 0 {
		return []model.Model{}, nil
	}
	query := r.db.WithContext(ctx).Model(&model.Model{})
	for i, pair := range pairs {
		if i == 0 {
			query = query.Where("name = ? AND provider = ?", pair.Name, pair.Provider)
			continue
		}
		query = query.Or("name = ? AND provider = ?", pair.Name, pair.Provider)
	}
	var models []model.Model
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *ModelRepo) ListByIDs(ctx context.Context, ids []string) ([]model.Model, error) {
	if len(ids) == 0 {
		return []model.Model{}, nil
	}
	var models []model.Model
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *ModelRepo) ListAllByProvider(ctx context.Context, provider string) ([]model.Model, error) {
	var models []model.Model
	if err := r.db.WithContext(ctx).
		Where("provider = ?", provider).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *ModelRepo) GetByID(ctx context.Context, id string) (*model.Model, error) {
	var item model.Model
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *ModelRepo) GetByName(ctx context.Context, name string) (*model.Model, error) {
	var item model.Model
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *ModelRepo) GetByNameProvider(ctx context.Context, name, provider string) (*model.Model, error) {
	var item model.Model
	err := r.db.WithContext(ctx).
		Where("name = ? AND provider = ?", name, provider).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *ModelRepo) Create(ctx context.Context, item *model.Model) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ModelRepo) Update(ctx context.Context, id string, updates map[string]any) (*model.Model, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Model{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *ModelRepo) UpsertBatch(ctx context.Context, items []model.Model) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "name"},
				{Name: "provider"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"status":     "active",
				"updated_at": gorm.Expr("now()"),
			}),
		}).
		Create(&items).Error
}

func (r *ModelRepo) ListProviders(ctx context.Context, activeOnly bool) ([]string, error) {
	query := r.db.WithContext(ctx).
		Model(&model.Model{}).
		Distinct("provider")
	if activeOnly {
		query = query.Where("status = ?", "active")
	}
	var providers []string
	if err := query.Order("provider ASC").Pluck("provider", &providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}
