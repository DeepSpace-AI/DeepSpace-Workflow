package newapi

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	BaseURL string
	APIKey  string
}

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
