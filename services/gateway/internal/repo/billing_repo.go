package repo

import (
	"context"
	"errors"
	"time"

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

func (r *BillingRepo) GetWalletForUpdate(ctx context.Context, userID int64) (*model.Wallet, error) {
	var w model.Wallet
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&w).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *BillingRepo) GetWallet(ctx context.Context, userID int64) (*model.Wallet, error) {
	var w model.Wallet
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&w).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *BillingRepo) CreateWallet(ctx context.Context, userID int64) (*model.Wallet, error) {
	wallet := model.Wallet{UserID: userID, Balance: 0, FrozenBalance: 0}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&wallet).Error
	if err != nil {
		return nil, err
	}
	return r.GetWalletForUpdate(ctx, userID)
}

func (r *BillingRepo) UpdateWallet(ctx context.Context, userID int64, balance, frozen float64) error {
	return r.db.WithContext(ctx).
		Model(&model.Wallet{}).
		Where("user_id = ?", userID).
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

func (r *BillingRepo) CreateTransaction(ctx context.Context, userID int64, typ string, amount float64, refID string, metadata []byte) (*model.Transaction, error) {
	tr := model.Transaction{
		UserID:   userID,
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

type WalletWithUser struct {
	UserID        int64     `json:"user_id"`
	Balance       float64   `json:"balance"`
	FrozenBalance float64   `json:"frozen_balance"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Status        string    `json:"status"`
	Role          string    `json:"role"`
	UserCreatedAt time.Time `json:"user_created_at"`
}

type WalletListFilter struct {
	UserID *int64
	Limit  int
	Offset int
}

func (r *BillingRepo) ListWallets(ctx context.Context, filter WalletListFilter) ([]WalletWithUser, int64, error) {
	query := r.db.WithContext(ctx).
		Table("wallets").
		Select("wallets.user_id, wallets.balance, wallets.frozen_balance, wallets.updated_at, users.email, users.status, users.role, users.created_at AS user_created_at").
		Joins("JOIN users ON users.id = wallets.user_id")
	if filter.UserID != nil {
		query = query.Where("wallets.user_id = ?", *filter.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []WalletWithUser
	if err := query.Order("wallets.updated_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

type TransactionListFilter struct {
	UserID *int64
	Type   string
	Start  *time.Time
	End    *time.Time
	Limit  int
	Offset int
}

func (r *BillingRepo) ListTransactions(ctx context.Context, filter TransactionListFilter) ([]model.Transaction, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Transaction{})
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Start != nil {
		query = query.Where("created_at >= ?", *filter.Start)
	}
	if filter.End != nil {
		query = query.Where("created_at < ?", *filter.End)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Transaction
	if err := query.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
