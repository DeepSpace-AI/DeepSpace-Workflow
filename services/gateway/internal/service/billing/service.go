package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/repo"

	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInsufficientFrozen  = errors.New("insufficient frozen balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrRefConflict         = errors.New("ref_id already used with different transaction")
)

type Service struct {
	db   *gorm.DB
	repo *repo.BillingRepo
}

func New(db *gorm.DB, repo *repo.BillingRepo) *Service {
	return &Service{db: db, repo: repo}
}

type HoldResult struct {
	Wallet      *model.Wallet      `json:"wallet"`
	Transaction *model.Transaction `json:"transaction"`
}

type CaptureResult = HoldResult

type ReleaseResult = HoldResult

type TopUpResult = HoldResult

type WalletListInput struct {
	UserID   *int64
	Page     int
	PageSize int
}

type TransactionListInput struct {
	UserID   *int64
	Type     string
	Start    *time.Time
	End      *time.Time
	Page     int
	PageSize int
}

func (s *Service) GetWallet(ctx context.Context, userID int64) (*model.Wallet, error) {
	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet != nil {
		return wallet, nil
	}
	return s.repo.CreateWallet(ctx, userID)
}

func (s *Service) Hold(ctx context.Context, userID int64, amount float64, refID string, metadata map[string]any) (*HoldResult, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	return s.withTx(ctx, func(repoTx *repo.BillingRepo) (*HoldResult, error) {
		if existing, err := s.findExisting(ctx, repoTx, userID, refID, "hold", amount); err != nil {
			return nil, err
		} else if existing != nil {
			wallet, err := s.ensureWallet(ctx, repoTx, userID)
			if err != nil {
				return nil, err
			}
			return &HoldResult{Wallet: wallet, Transaction: existing}, nil
		}

		wallet, err := s.ensureWallet(ctx, repoTx, userID)
		if err != nil {
			return nil, err
		}

		if wallet.Balance < amount {
			return nil, ErrInsufficientBalance
		}

		wallet.Balance -= amount
		wallet.FrozenBalance += amount
		if err := repoTx.UpdateWallet(ctx, userID, wallet.Balance, wallet.FrozenBalance); err != nil {
			return nil, err
		}

		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		tr, err := repoTx.CreateTransaction(ctx, userID, "hold", amount, refID, meta)
		if err != nil {
			return nil, err
		}

		return &HoldResult{Wallet: wallet, Transaction: tr}, nil
	})
}

func (s *Service) Capture(ctx context.Context, userID int64, amount float64, refID string, metadata map[string]any) (*CaptureResult, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	return s.withTx(ctx, func(repoTx *repo.BillingRepo) (*CaptureResult, error) {
		if existing, err := s.findExisting(ctx, repoTx, userID, refID, "capture", amount); err != nil {
			return nil, err
		} else if existing != nil {
			wallet, err := s.ensureWallet(ctx, repoTx, userID)
			if err != nil {
				return nil, err
			}
			return &HoldResult{Wallet: wallet, Transaction: existing}, nil
		}

		wallet, err := s.ensureWallet(ctx, repoTx, userID)
		if err != nil {
			return nil, err
		}

		if wallet.FrozenBalance < amount {
			return nil, ErrInsufficientFrozen
		}

		wallet.FrozenBalance -= amount
		if err := repoTx.UpdateWallet(ctx, userID, wallet.Balance, wallet.FrozenBalance); err != nil {
			return nil, err
		}

		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		tr, err := repoTx.CreateTransaction(ctx, userID, "capture", amount, refID, meta)
		if err != nil {
			return nil, err
		}

		return &HoldResult{Wallet: wallet, Transaction: tr}, nil
	})
}

func (s *Service) TopUp(ctx context.Context, userID int64, amount float64, currency string, refID string, metadata map[string]any) (*TopUpResult, error) {
	if amount == 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return nil, ErrInvalidAmount
	}
	return s.withTx(ctx, func(repoTx *repo.BillingRepo) (*HoldResult, error) {
		if existing, err := s.findExisting(ctx, repoTx, userID, refID, "topup", amount); err != nil {
			return nil, err
		} else if existing != nil {
			wallet, err := s.ensureWallet(ctx, repoTx, userID)
			if err != nil {
				return nil, err
			}
			return &HoldResult{Wallet: wallet, Transaction: existing}, nil
		}

		wallet, err := s.ensureWallet(ctx, repoTx, userID)
		if err != nil {
			return nil, err
		}

		wallet.Balance += amount
		if err := repoTx.UpdateWallet(ctx, userID, wallet.Balance, wallet.FrozenBalance); err != nil {
			return nil, err
		}

		if metadata == nil {
			metadata = map[string]any{}
		}
		metadata["currency"] = normalizeCurrency(currency)

		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		tr, err := repoTx.CreateTransaction(ctx, userID, "topup", amount, refID, meta)
		if err != nil {
			return nil, err
		}

		return &HoldResult{Wallet: wallet, Transaction: tr}, nil
	})
}

func (s *Service) Release(ctx context.Context, userID int64, amount float64, refID string, metadata map[string]any) (*ReleaseResult, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	return s.withTx(ctx, func(repoTx *repo.BillingRepo) (*ReleaseResult, error) {
		if existing, err := s.findExisting(ctx, repoTx, userID, refID, "release", amount); err != nil {
			return nil, err
		} else if existing != nil {
			wallet, err := s.ensureWallet(ctx, repoTx, userID)
			if err != nil {
				return nil, err
			}
			return &HoldResult{Wallet: wallet, Transaction: existing}, nil
		}

		wallet, err := s.ensureWallet(ctx, repoTx, userID)
		if err != nil {
			return nil, err
		}

		if wallet.FrozenBalance < amount {
			return nil, ErrInsufficientFrozen
		}

		wallet.FrozenBalance -= amount
		wallet.Balance += amount
		if err := repoTx.UpdateWallet(ctx, userID, wallet.Balance, wallet.FrozenBalance); err != nil {
			return nil, err
		}

		meta, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		tr, err := repoTx.CreateTransaction(ctx, userID, "release", amount, refID, meta)
		if err != nil {
			return nil, err
		}

		return &HoldResult{Wallet: wallet, Transaction: tr}, nil
	})
}

func (s *Service) ListWallets(ctx context.Context, input WalletListInput) ([]repo.WalletWithUser, int64, error) {
	page, pageSize := normalizePage(input.Page, input.PageSize)
	return s.repo.ListWallets(ctx, repo.WalletListFilter{
		UserID: input.UserID,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
}

func (s *Service) ListTransactions(ctx context.Context, input TransactionListInput) ([]model.Transaction, int64, error) {
	page, pageSize := normalizePage(input.Page, input.PageSize)
	return s.repo.ListTransactions(ctx, repo.TransactionListFilter{
		UserID: input.UserID,
		Type:   input.Type,
		Start:  input.Start,
		End:    input.End,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
}

func (s *Service) withTx(ctx context.Context, fn func(repoTx *repo.BillingRepo) (*HoldResult, error)) (*HoldResult, error) {
	var result *HoldResult
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.repo.WithTx(tx)
		value, err := fn(repoTx)
		if err != nil {
			return err
		}
		result = value
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) ensureWallet(ctx context.Context, repoTx *repo.BillingRepo, userID int64) (*model.Wallet, error) {
	wallet, err := repoTx.GetWalletForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		wallet, err = repoTx.CreateWallet(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (s *Service) findExisting(ctx context.Context, repoTx *repo.BillingRepo, userID int64, refID, typ string, amount float64) (*model.Transaction, error) {
	if refID == "" {
		return nil, fmt.Errorf("ref_id is required")
	}
	existing, err := repoTx.GetTransactionByRef(ctx, refID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	if existing.UserID != userID {
		return nil, ErrRefConflict
	}
	if existing.Type == typ {
		if existing.Amount != amount {
			return nil, ErrRefConflict
		}
		return existing, nil
	}
	return nil, ErrRefConflict
}

func normalizeCurrency(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "CNY"
	}
	return strings.ToUpper(value)
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
