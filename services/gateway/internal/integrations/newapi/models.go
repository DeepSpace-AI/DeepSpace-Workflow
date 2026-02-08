package newapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ModelItem struct {
	ID string `json:"id"`
}

type UpstreamModel struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type listModelsResponse struct {
	Data []ModelItem `json:"data"`
}

func (c *Client) ListModels(ctx context.Context) ([]UpstreamModel, error) {
	if c == nil {
		return nil, fmt.Errorf("client missing")
	}
	base := strings.TrimRight(c.BaseURL, "/")
	target := base + "/v1/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload listModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	items := make([]UpstreamModel, 0, len(payload.Data))
	for _, item := range payload.Data {
		name := strings.TrimSpace(item.ID)
		if name == "" {
			continue
		}
		items = append(items, UpstreamModel{
			Name:     name,
			Provider: guessProvider(name),
		})
	}
	return items, nil
}

func guessProvider(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "gpt"):
		return "openai"
	case strings.Contains(lower, "claude"):
		return "anthropic"
	case strings.Contains(lower, "deepseek"):
		return "deepseek"
	case strings.Contains(lower, "qwen"):
		return "qwen"
	default:
		return "other"
	}
}
