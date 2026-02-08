package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type KnowledgeBaseWithCount struct {
	model.KnowledgeBase
	DocCount int64 `gorm:"column:doc_count"`
}

type KnowledgeRepo struct {
	db *gorm.DB
}

func NewKnowledgeRepo(db *gorm.DB) *KnowledgeRepo {
	return &KnowledgeRepo{db: db}
}

func (r *KnowledgeRepo) ListBases(ctx context.Context, orgID int64, scope string, projectID *int64, includeOrg bool) ([]KnowledgeBaseWithCount, error) {
	query := r.db.WithContext(ctx).
		Table("knowledge_bases").
		Select("knowledge_bases.*, COUNT(knowledge_documents.id) AS doc_count").
		Joins("LEFT JOIN knowledge_documents ON knowledge_documents.knowledge_base_id = knowledge_bases.id").
		Where("knowledge_bases.user_id = ?", orgID).
		Group("knowledge_bases.id").
		Order("knowledge_bases.created_at DESC")

	switch scope {
	case "org":
		query = query.Where("knowledge_bases.scope = ?", "org")
	case "project":
		if projectID == nil {
			return []KnowledgeBaseWithCount{}, nil
		}
		if includeOrg {
			query = query.Where("(knowledge_bases.scope = 'org' OR (knowledge_bases.scope = 'project' AND knowledge_bases.project_id = ?))", *projectID)
		} else {
			query = query.Where("knowledge_bases.scope = 'project' AND knowledge_bases.project_id = ?", *projectID)
		}
	case "all":
		if projectID != nil {
			query = query.Where("(knowledge_bases.scope = 'org' OR (knowledge_bases.scope = 'project' AND knowledge_bases.project_id = ?))", *projectID)
		}
	}

	var bases []KnowledgeBaseWithCount
	if err := query.Find(&bases).Error; err != nil {
		return nil, err
	}
	return bases, nil
}

func (r *KnowledgeRepo) GetBase(ctx context.Context, orgID, id int64) (*model.KnowledgeBase, error) {
	var base model.KnowledgeBase
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, orgID).
		First(&base).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &base, nil
}

func (r *KnowledgeRepo) CreateBase(ctx context.Context, base *model.KnowledgeBase) error {
	return r.db.WithContext(ctx).Create(base).Error
}

func (r *KnowledgeRepo) UpdateBase(ctx context.Context, orgID, id int64, updates map[string]any) (*model.KnowledgeBase, error) {
	if len(updates) == 0 {
		return r.GetBase(ctx, orgID, id)
	}
	err := r.db.WithContext(ctx).
		Model(&model.KnowledgeBase{}).
		Where("id = ? AND user_id = ?", id, orgID).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return r.GetBase(ctx, orgID, id)
}

func (r *KnowledgeRepo) DeleteBase(ctx context.Context, orgID, id int64) (bool, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, orgID).
		Delete(&model.KnowledgeBase{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *KnowledgeRepo) ListDocuments(ctx context.Context, orgID, kbID int64) ([]model.KnowledgeDocument, error) {
	var docs []model.KnowledgeDocument
	if err := r.db.WithContext(ctx).
		Where("knowledge_base_id = ? AND user_id = ?", kbID, orgID).
		Order("created_at DESC").
		Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *KnowledgeRepo) CreateDocument(ctx context.Context, doc *model.KnowledgeDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *KnowledgeRepo) DeleteDocument(ctx context.Context, orgID, kbID, docID int64) (*model.KnowledgeDocument, error) {
	var doc model.KnowledgeDocument
	err := r.db.WithContext(ctx).
		Where("id = ? AND knowledge_base_id = ? AND user_id = ?", docID, kbID, orgID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.WithContext(ctx).Delete(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *KnowledgeRepo) GetDocument(ctx context.Context, orgID, kbID, docID int64) (*model.KnowledgeDocument, error) {
	var doc model.KnowledgeDocument
	err := r.db.WithContext(ctx).
		Where("id = ? AND knowledge_base_id = ? AND user_id = ?", docID, kbID, orgID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *KnowledgeRepo) CountDocumentsByOrg(ctx context.Context, orgID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.KnowledgeDocument{}).
		Where("user_id = ?", orgID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
