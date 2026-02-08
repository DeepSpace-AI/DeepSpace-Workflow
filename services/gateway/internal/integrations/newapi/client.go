package newapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	BaseURL string
	APIKey  string
}

type ParsedUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

const usageContextKey = "newapi_usage"

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

// ProxyRequestHandler creates a reverse proxy to NewAPI
func (c *Client) ProxyRequestHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.Proxy(ctx)
	}
}

// Proxy forwards the current request to NewAPI.
func (c *Client) Proxy(ctx *gin.Context) {
	remote, err := url.Parse(c.BaseURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid NewAPI Base URL"})
		ctx.Abort()
		return
	}

	// Improve streaming behavior for SSE responses.
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.FlushInterval = 100 * time.Millisecond
	proxy.ModifyResponse = func(resp *http.Response) error {
		contentType := strings.ToLower(resp.Header.Get("Content-Type"))
		isSSE := strings.Contains(contentType, "text/event-stream")
		resp.Body = &usageCaptureReadCloser{
			rc:    resp.Body,
			isSSE: isSSE,
			onUsage: func(usage ParsedUsage) {
				ctx.Set(usageContextKey, usage)
			},
		}
		return nil
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host

		// If NewAPIKey is configured, use it.
		// Otherwise, we might rely on the incoming request's key (if valid).
		// But typically Gateway -> NewAPI uses a high-privilege key or a specific token.
		if c.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.APIKey)
		}
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func GetUsageFromContext(ctx *gin.Context) (ParsedUsage, bool) {
	if ctx == nil {
		return ParsedUsage{}, false
	}
	value, ok := ctx.Get(usageContextKey)
	if !ok {
		return ParsedUsage{}, false
	}
	usage, ok := value.(ParsedUsage)
	return usage, ok
}

type usageCaptureReadCloser struct {
	rc       io.ReadCloser
	isSSE    bool
	onUsage  func(ParsedUsage)
	usageSet bool
	buf      []byte
	jsonBuf  bytes.Buffer
}

func (u *usageCaptureReadCloser) Read(p []byte) (int, error) {
	n, err := u.rc.Read(p)
	if n > 0 {
		if u.isSSE {
			u.consumeSSE(p[:n])
		} else {
			u.jsonBuf.Write(p[:n])
		}
	}
	if err == io.EOF {
		u.finalize()
	}
	return n, err
}

func (u *usageCaptureReadCloser) Close() error {
	u.finalize()
	return u.rc.Close()
}

func (u *usageCaptureReadCloser) consumeSSE(chunk []byte) {
	if u.usageSet {
		return
	}
	u.buf = append(u.buf, chunk...)
	for {
		idx := bytes.IndexByte(u.buf, '\n')
		if idx < 0 {
			break
		}
		line := strings.TrimSpace(strings.TrimRight(string(u.buf[:idx]), "\r"))
		u.buf = u.buf[idx+1:]
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			continue
		}
		usage, ok := parseUsageFromJSON([]byte(payload))
		if ok {
			u.setUsage(usage)
			return
		}
	}
}

func (u *usageCaptureReadCloser) finalize() {
	if u.usageSet {
		return
	}
	if u.isSSE {
		if len(u.buf) == 0 {
			return
		}
		line := strings.TrimSpace(strings.TrimRight(string(u.buf), "\r"))
		if strings.HasPrefix(line, "data:") {
			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload != "" && payload != "[DONE]" {
				if usage, ok := parseUsageFromJSON([]byte(payload)); ok {
					u.setUsage(usage)
				}
			}
		}
		return
	}
	if u.jsonBuf.Len() == 0 {
		return
	}
	if usage, ok := parseUsageFromJSON(u.jsonBuf.Bytes()); ok {
		u.setUsage(usage)
	}
}

func (u *usageCaptureReadCloser) setUsage(usage ParsedUsage) {
	if u.usageSet {
		return
	}
	u.usageSet = true
	if u.onUsage != nil {
		u.onUsage(usage)
	}
}

func parseUsageFromJSON(data []byte) (ParsedUsage, bool) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return ParsedUsage{}, false
	}
	if usageRaw, ok := raw["usage"].(map[string]any); ok {
		return parseUsageObject(usageRaw), true
	}
	if responseRaw, ok := raw["response"].(map[string]any); ok {
		if usageRaw, ok := responseRaw["usage"].(map[string]any); ok {
			return parseUsageObject(usageRaw), true
		}
	}
	if usageRaw, ok := raw["usageMetadata"].(map[string]any); ok {
		return parseGeminiUsage(usageRaw), true
	}
	return ParsedUsage{}, false
}

func parseUsageObject(usage map[string]any) ParsedUsage {
	prompt := getIntValue(usage["prompt_tokens"])
	completion := getIntValue(usage["completion_tokens"])
	total := getIntValue(usage["total_tokens"])
	if prompt == 0 && completion == 0 {
		prompt = getIntValue(usage["input_tokens"])
		completion = getIntValue(usage["output_tokens"])
	}
	if total == 0 && (prompt > 0 || completion > 0) {
		total = prompt + completion
	}
	return ParsedUsage{PromptTokens: prompt, CompletionTokens: completion, TotalTokens: total}
}

func parseGeminiUsage(usage map[string]any) ParsedUsage {
	prompt := getIntValue(usage["promptTokenCount"])
	completion := getIntValue(usage["candidatesTokenCount"])
	total := getIntValue(usage["totalTokenCount"])
	if total == 0 && (prompt > 0 || completion > 0) {
		total = prompt + completion
	}
	return ParsedUsage{PromptTokens: prompt, CompletionTokens: completion, TotalTokens: total}
}

func getIntValue(value any) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			return 0
		}
		return int(parsed)
	default:
		return 0
	}
}
