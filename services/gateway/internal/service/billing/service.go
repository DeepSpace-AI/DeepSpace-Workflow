package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
