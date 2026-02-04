package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BillingRepo struct {
	db *gorm.DB
}

func NewBillingRepo(db *gorm.DB) *BillingRepo {
	return &BillingRepo{db: db}
}

func (r *BillingRepo) WithTx(tx *gorm.DB) *BillingRepo {
	return &BillingRepo{db: tx}
}

func (r *BillingRepo) GetWalletForUpdate(ctx context.Context, orgID int64) (*model.Wallet, error) {
	var w model.Wallet
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("org_id = ?", orgID).
		First(&w).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *BillingRepo) CreateWallet(ctx context.Context, orgID int64) (*model.Wallet, error) {
	wallet := model.Wallet{OrgID: orgID, Balance: 0, FrozenBalance: 0}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&wallet).Error
	if err != nil {
		return nil, err
	}
	return r.GetWalletForUpdate(ctx, orgID)
}

func (r *BillingRepo) UpdateWallet(ctx context.Context, orgID int64, balance, frozen float64) error {
	return r.db.WithContext(ctx).
		Model(&model.Wallet{}).
		Where("org_id = ?", orgID).
		Updates(map[string]any{
			"balance":        balance,
			"frozen_balance": frozen,
		}).Error
}

func (r *BillingRepo) GetTransactionByRef(ctx context.Context, refID string) (*model.Transaction, error) {
	var t model.Transaction
	err := r.db.WithContext(ctx).
		Where("ref_id = ?", refID).
		First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *BillingRepo) CreateTransaction(ctx context.Context, orgID int64, typ string, amount float64, refID string, metadata []byte) (*model.Transaction, error) {
	tr := model.Transaction{
		OrgID:    orgID,
		Type:     typ,
		Amount:   amount,
		RefID:    refID,
		Metadata: metadata,
	}

	if err := r.db.WithContext(ctx).Create(&tr).Error; err != nil {
		return nil, err
	}
	return &tr, nil
}
