package apikey

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"deepspace/internal/pkg/keys"
	"deepspace/internal/repo"
)

const defaultKeyPrefixLen = 8

type Service struct {
	repo *repo.APIKeyRepo
}

func New(repo *repo.APIKeyRepo) *Service {
	return &Service{repo: repo}
}

type CreatedKey struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Prefix    string    `json:"prefix"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Key       string    `json:"key"`
	Scopes    []string  `json:"scopes"`
}

type KeyItem struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	Scopes     []string   `json:"scopes"`
}

func (s *Service) Create(ctx context.Context, orgID int64, name string, scopes []string) (*CreatedKey, error) {
	plain := keys.GenerateAPIKey()
	hash := keys.HashAPIKey(plain)
	prefix := keys.KeyPrefix(plain, defaultKeyPrefixLen)

	scopesJSON, err := json.Marshal(scopes)
	if err != nil {
		return nil, err
	}

	var namePtr *string
	if strings.TrimSpace(name) != "" {
		trimmed := strings.TrimSpace(name)
		namePtr = &trimmed
	}

	record, err := s.repo.Create(ctx, orgID, namePtr, hash, prefix, scopesJSON)
	if err != nil {
		return nil, err
	}

	recordName := ""
	if record.Name != nil {
		recordName = *record.Name
	}

	return &CreatedKey{
		ID:        record.ID,
		Name:      recordName,
		Prefix:    record.KeyPrefix,
		Status:    record.Status,
		CreatedAt: record.CreatedAt,
		Key:       plain,
		Scopes:    scopes,
	}, nil
}

func (s *Service) List(ctx context.Context, orgID int64) ([]KeyItem, error) {
	items, err := s.repo.ListByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	result := make([]KeyItem, 0, len(items))
	for _, item := range items {
		var lastUsed *time.Time
		if item.LastUsedAt != nil {
			lastUsed = item.LastUsedAt
		}
		name := ""
		if item.Name != nil {
			name = *item.Name
		}
		var scopes []string
		if len(item.Scopes) > 0 {
			_ = json.Unmarshal(item.Scopes, &scopes)
		}
		result = append(result, KeyItem{
			ID:         item.ID,
			Name:       name,
			Prefix:     item.KeyPrefix,
			Status:     item.Status,
			CreatedAt:  item.CreatedAt,
			LastUsedAt: lastUsed,
			Scopes:     scopes,
		})
	}

	return result, nil
}

func (s *Service) Disable(ctx context.Context, orgID, keyID int64) (bool, error) {
	return s.repo.Disable(ctx, orgID, keyID)
}

func (s *Service) Delete(ctx context.Context, orgID, keyID int64) (bool, error) {
	return s.repo.Delete(ctx, orgID, keyID)
}

func (s *Service) UpdateScopes(ctx context.Context, orgID, keyID int64, scopes []string) (bool, error) {
	scopesJSON, err := json.Marshal(scopes)
	if err != nil {
		return false, err
	}
	return s.repo.UpdateScopes(ctx, orgID, keyID, scopesJSON)
}
