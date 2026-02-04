package auth

import (
	"context"
	"errors"

	"deepspace/internal/pkg/keys"
	"deepspace/internal/repo"
)

var ErrInvalidAPIKey = errors.New("invalid api key")

// AuthContext carries resolved auth identity for the request.
type AuthContext struct {
	APIKeyID int64
	OrgID    int64
}

type APIKeyValidator struct {
	repo *repo.APIKeyRepo
}

func NewAPIKeyValidator(repo *repo.APIKeyRepo) *APIKeyValidator {
	return &APIKeyValidator{repo: repo}
}

func (v *APIKeyValidator) Validate(ctx context.Context, apiKey string) (*AuthContext, error) {
	keyHash := keys.HashAPIKey(apiKey)
	key, err := v.repo.FindActiveByHash(ctx, keyHash)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, ErrInvalidAPIKey
	}

	_ = v.repo.UpdateLastUsed(ctx, key.ID)

	return &AuthContext{
		APIKeyID: key.ID,
		OrgID:    key.OrgID,
	}, nil
}
