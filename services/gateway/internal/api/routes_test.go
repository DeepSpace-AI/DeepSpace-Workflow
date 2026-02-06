package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"deepspace/internal/config"
	"deepspace/internal/service/auth"
	"deepspace/internal/service/billing"
	"deepspace/internal/service/chat"
	"deepspace/internal/service/knowledge"
	"deepspace/internal/service/project"
	"deepspace/internal/service/projectdocument"
	"deepspace/internal/service/projectskill"
	"deepspace/internal/service/projectworkflow"
	"deepspace/internal/service/usage"
	"deepspace/internal/service/user"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := &config.Config{
		Port:          "8080",
		NewAPIBaseURL: "http://example.com",
		NewAPIKey:     "test-key",
	}

	SetupRoutes(
		r,
		cfg,
		(*billing.Service)(nil),
		(*usage.Service)(nil),
		(*project.Service)(nil),
		(*chat.Service)(nil),
		(*knowledge.Service)(nil),
		(*projectdocument.Service)(nil),
		(*projectskill.Service)(nil),
		(*projectworkflow.Service)(nil),
		(*auth.UserAuthService)(nil),
		(*user.Service)(nil),
		(*auth.JWTManager)(nil),
	)

	// Test Health
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	expected := `{"status":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("Expected %s, got %s", expected, w.Body.String())
	}
}
