package model

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"

	"github.com/google/uuid"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

var (
	ErrInvalidName         = errors.New("invalid name")
	ErrInvalidProvider     = errors.New("invalid provider")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidCapability   = errors.New("invalid capability")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrDuplicateName       = errors.New("duplicate name")
	ErrModelNotFound       = errors.New("model not found")
	ErrInvalidCurrency     = errors.New("invalid currency")
	ErrInvalidCapabilities = errors.New("invalid capabilities")
)

var allowedCapabilities = map[string]struct{}{
	"chat":          {},
	"completion":    {},
	"embedding":     {},
	"vision":        {},
	"image":         {},
	"audio":         {},
	"tool":          {},
	"skill":         {},
	"stream":        {},
	"json_mode":     {},
	"function_call": {},
}

type Service struct {
	repo *repo.ModelRepo
}

func New(repo *repo.ModelRepo) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	Name         string
	Provider     string
	PriceInput   float64
	PriceOutput  float64
	Currency     string
	Capabilities []string
	Status       string
	Metadata     map[string]any
}

type UpdateInput struct {
	Provider     *string
	PriceInput   *float64
	PriceOutput  *float64
	Currency     *string
	Capabilities *[]string
	Status       *string
	ProviderIcon *string
	Metadata     *map[string]any
}

type ModelItem struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Provider     string         `json:"provider"`
	PriceInput   float64        `json:"price_input"`
	PriceOutput  float64        `json:"price_output"`
	Currency     string         `json:"currency"`
	Capabilities []string       `json:"capabilities"`
	Status       string         `json:"status"`
	Metadata     map[string]any `json:"metadata"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

type ConfirmInput struct {
	Name     string
	Provider string
}

type ConfirmResult struct {
	Items   []ModelItem `json:"items"`
	Created int         `json:"created"`
	Updated int         `json:"updated"`
}

type BatchPricingItem struct {
	ID           string
	PriceInput   *float64
	PriceOutput  *float64
	Currency     *string
	Status       *string
	Capabilities *[]string
}

type BatchPricingResult struct {
	Items     []ModelItem `json:"items"`
	MissedIDs []string    `json:"missed_ids"`
}

func (s *Service) ListActive(ctx context.Context) ([]ModelItem, error) {
	items, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]ModelItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapModelItem(&item))
	}
	return result, nil
}

func (s *Service) ListActiveByProvider(ctx context.Context, provider string) ([]ModelItem, error) {
	items, err := s.repo.ListActiveByProvider(ctx, provider)
	if err != nil {
		return nil, err
	}
	result := make([]ModelItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapModelItem(&item))
	}
	return result, nil
}

func (s *Service) ListAll(ctx context.Context) ([]ModelItem, error) {
	items, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]ModelItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapModelItem(&item))
	}
	return result, nil
}

func (s *Service) ListAllByProvider(ctx context.Context, provider string) ([]ModelItem, error) {
	items, err := s.repo.ListAllByProvider(ctx, provider)
	if err != nil {
		return nil, err
	}
	result := make([]ModelItem, 0, len(items))
	for _, item := range items {
		result = append(result, mapModelItem(&item))
	}
	return result, nil
}

func (s *Service) ListProviders(ctx context.Context, activeOnly bool) ([]string, error) {
	return s.repo.ListProviders(ctx, activeOnly)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*ModelItem, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidName
	}
	provider := strings.TrimSpace(input.Provider)
	if provider == "" {
		return nil, ErrInvalidProvider
	}
	if input.PriceInput < 0 || input.PriceOutput < 0 {
		return nil, ErrInvalidPrice
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		return nil, ErrInvalidStatus
	}
	currency := normalizeCurrency(input.Currency)
	if currency == "" {
		return nil, ErrInvalidCurrency
	}

	capabilities, err := normalizeCapabilities(input.Capabilities)
	if err != nil {
		return nil, err
	}

	if existing, err := s.repo.GetByNameProvider(ctx, name, provider); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, ErrDuplicateName
	}

	capsJSON, _ := json.Marshal(capabilities)
	metadataJSON, _ := json.Marshal(normalizeMetadata(input.Metadata))

	item := &model.Model{
		ID:           uuid.New(),
		Name:         name,
		Provider:     provider,
		PriceInput:   input.PriceInput,
		PriceOutput:  input.PriceOutput,
		Currency:     currency,
		Capabilities: capsJSON,
		Status:       status,
		Metadata:     metadataJSON,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	mapped := mapModelItem(item)
	return &mapped, nil
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*ModelItem, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrModelNotFound
	}
	updates := map[string]any{}

	if input.Provider != nil {
		value := strings.TrimSpace(*input.Provider)
		if value == "" {
			return nil, ErrInvalidProvider
		}
		updates["provider"] = value
	}
	if input.PriceInput != nil {
		if *input.PriceInput < 0 {
			return nil, ErrInvalidPrice
		}
		updates["price_input"] = *input.PriceInput
	}
	if input.PriceOutput != nil {
		if *input.PriceOutput < 0 {
			return nil, ErrInvalidPrice
		}
		updates["price_output"] = *input.PriceOutput
	}
	if input.Currency != nil {
		value := normalizeCurrency(*input.Currency)
		if value == "" {
			return nil, ErrInvalidCurrency
		}
		updates["currency"] = value
	}
	if input.Capabilities != nil {
		caps, err := normalizeCapabilities(*input.Capabilities)
		if err != nil {
			return nil, err
		}
		capsJSON, _ := json.Marshal(caps)
		updates["capabilities"] = capsJSON
	}
	if input.Status != nil {
		value := normalizeStatus(*input.Status)
		if value == "" {
			return nil, ErrInvalidStatus
		}
		updates["status"] = value
	}
	metadataUpdated := false
	metadata := map[string]any{}
	if input.Metadata != nil {
		metadata = normalizeMetadata(*input.Metadata)
		metadataUpdated = true
	}
	if input.ProviderIcon != nil {
		if !metadataUpdated {
			existing, err := s.repo.GetByID(ctx, id)
			if err != nil {
				return nil, err
			}
			if existing != nil && len(existing.Metadata) > 0 {
				_ = json.Unmarshal(existing.Metadata, &metadata)
			}
			metadataUpdated = true
		}
		metadata["provider_icon"] = *input.ProviderIcon
	}
	if metadataUpdated {
		metadataJSON, _ := json.Marshal(normalizeMetadata(metadata))
		updates["metadata"] = metadataJSON
	}

	item, err := s.repo.Update(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrModelNotFound
	}
	mapped := mapModelItem(item)
	return &mapped, nil
}

func (s *Service) ConfirmBatch(ctx context.Context, inputs []ConfirmInput) (*ConfirmResult, error) {
	if len(inputs) == 0 {
		return &ConfirmResult{Items: []ModelItem{}}, nil
	}

	unique := map[string]ConfirmInput{}
	for _, input := range inputs {
		name := strings.TrimSpace(input.Name)
		if name == "" {
			return nil, ErrInvalidName
		}
		provider := strings.TrimSpace(input.Provider)
		if provider == "" {
			return nil, ErrInvalidProvider
		}
		key := buildModelKey(name, provider)
		unique[key] = ConfirmInput{Name: name, Provider: provider}
	}

	pairs := make([]repo.NameProvider, 0, len(unique))
	items := make([]model.Model, 0, len(unique))
	for _, input := range unique {
		pairs = append(pairs, repo.NameProvider{Name: input.Name, Provider: input.Provider})
		capsJSON, _ := json.Marshal([]string{})
		metadataJSON, _ := json.Marshal(map[string]any{})
		items = append(items, model.Model{
			ID:           uuid.New(),
			Name:         input.Name,
			Provider:     input.Provider,
			PriceInput:   0,
			PriceOutput:  0,
			Currency:     "USD",
			Capabilities: capsJSON,
			Status:       StatusActive,
			Metadata:     metadataJSON,
		})
	}

	existing, err := s.repo.ListByPairs(ctx, pairs)
	if err != nil {
		return nil, err
	}
	existingKeys := map[string]struct{}{}
	for _, item := range existing {
		key := buildModelKey(item.Name, item.Provider)
		existingKeys[key] = struct{}{}
	}

	if err := s.repo.UpsertBatch(ctx, items); err != nil {
		return nil, err
	}

	resultItems, err := s.repo.ListByPairs(ctx, pairs)
	if err != nil {
		return nil, err
	}

	result := &ConfirmResult{
		Items:   make([]ModelItem, 0, len(resultItems)),
		Created: len(pairs) - len(existingKeys),
		Updated: len(existingKeys),
	}
	for _, item := range resultItems {
		result.Items = append(result.Items, mapModelItem(&item))
	}
	return result, nil
}

func (s *Service) BatchPricing(ctx context.Context, items []BatchPricingItem) (*BatchPricingResult, error) {
	if len(items) == 0 {
		return &BatchPricingResult{Items: []ModelItem{}, MissedIDs: []string{}}, nil
	}

	uniqueIDs := map[string]struct{}{}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			continue
		}
		if _, exists := uniqueIDs[id]; exists {
			continue
		}
		uniqueIDs[id] = struct{}{}
		ids = append(ids, id)
	}

	current, err := s.repo.ListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	currentMap := map[string]*model.Model{}
	for i := range current {
		item := current[i]
		currentMap[item.ID.String()] = &item
	}

	result := &BatchPricingResult{
		Items:     make([]ModelItem, 0, len(items)),
		MissedIDs: make([]string, 0),
	}

	for _, input := range items {
		id := strings.TrimSpace(input.ID)
		if id == "" {
			result.MissedIDs = append(result.MissedIDs, id)
			continue
		}
		currentItem, ok := currentMap[id]
		if !ok || currentItem == nil {
			result.MissedIDs = append(result.MissedIDs, id)
			continue
		}

		updates := map[string]any{}
		if input.PriceInput != nil {
			if *input.PriceInput < 0 {
				return nil, ErrInvalidPrice
			}
			updates["price_input"] = *input.PriceInput
		}
		if input.PriceOutput != nil {
			if *input.PriceOutput < 0 {
				return nil, ErrInvalidPrice
			}
			updates["price_output"] = *input.PriceOutput
		}
		if input.Currency != nil {
			value := normalizeCurrency(*input.Currency)
			if value == "" {
				return nil, ErrInvalidCurrency
			}
			updates["currency"] = value
		}
		if input.Status != nil {
			value := normalizeStatus(*input.Status)
			if value == "" {
				return nil, ErrInvalidStatus
			}
			updates["status"] = value
		}
		if input.Capabilities != nil {
			caps, err := normalizeCapabilities(*input.Capabilities)
			if err != nil {
				return nil, err
			}
			capsJSON, _ := json.Marshal(caps)
			updates["capabilities"] = capsJSON
		}

		if len(updates) == 0 {
			result.Items = append(result.Items, mapModelItem(currentItem))
			continue
		}

		updated, err := s.repo.Update(ctx, id, updates)
		if err != nil {
			return nil, err
		}
		if updated == nil {
			result.MissedIDs = append(result.MissedIDs, id)
			continue
		}
		result.Items = append(result.Items, mapModelItem(updated))
	}

	return result, nil
}

func buildModelKey(name, provider string) string {
	return strings.ToLower(strings.TrimSpace(name)) + "|" + strings.ToLower(strings.TrimSpace(provider))
}

func mapModelItem(item *model.Model) ModelItem {
	capabilities := []string{}
	if len(item.Capabilities) > 0 {
		_ = json.Unmarshal(item.Capabilities, &capabilities)
	}
	metadata := map[string]any{}
	if len(item.Metadata) > 0 {
		_ = json.Unmarshal(item.Metadata, &metadata)
	}
	return ModelItem{
		ID:           item.ID.String(),
		Name:         item.Name,
		Provider:     item.Provider,
		PriceInput:   item.PriceInput,
		PriceOutput:  item.PriceOutput,
		Currency:     item.Currency,
		Capabilities: capabilities,
		Status:       item.Status,
		Metadata:     metadata,
		CreatedAt:    item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func normalizeStatus(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "":
		return StatusActive
	case StatusActive, StatusDisabled:
		return normalized
	default:
		return ""
	}
}

func normalizeCurrency(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "USD"
	}
	return strings.ToUpper(value)
}

func normalizeCapabilities(items []string) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}
	seen := map[string]struct{}{}
	result := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.ToLower(strings.TrimSpace(item))
		if value == "" {
			continue
		}
		if _, ok := allowedCapabilities[value]; !ok {
			return nil, ErrInvalidCapability
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result, nil
}

func normalizeMetadata(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	return input
}
