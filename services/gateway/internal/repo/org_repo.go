package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type OrgRepo struct {
	db *gorm.DB
}

func NewOrgRepo(db *gorm.DB) *OrgRepo {
	return &OrgRepo{db: db}
}

func (r *OrgRepo) Create(ctx context.Context, name string, ownerUserID int64) (*model.Org, error) {
	org := model.Org{Name: name, OwnerUserID: ownerUserID}
	if err := r.db.WithContext(ctx).Create(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrgRepo) GetByOwner(ctx context.Context, ownerUserID int64) (*model.Org, error) {
	var org model.Org
	err := r.db.WithContext(ctx).Where("owner_user_id = ?", ownerUserID).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &org, nil
}

func (r *OrgRepo) AddMember(ctx context.Context, orgID, userID int64, role string) error {
	member := model.OrgMember{OrgID: orgID, UserID: userID, Role: role}
	return r.db.WithContext(ctx).Create(&member).Error
}
