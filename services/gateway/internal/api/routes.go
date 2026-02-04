package api

import (
	"net/http"

	"deepspace/internal/api/handlers"
	"deepspace/internal/api/middleware"
	"deepspace/internal/config"
	"deepspace/internal/integrations/newapi"
	"deepspace/internal/service/apikey"
	"deepspace/internal/service/auth"
	"deepspace/internal/service/billing"
	"deepspace/internal/service/chat"
	"deepspace/internal/service/knowledge"
	"deepspace/internal/service/project"
	"deepspace/internal/service/usage"
	"deepspace/internal/service/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	apiKeyService *apikey.Service,
	billingService *billing.Service,
	usageService *usage.Service,
	projectService *project.Service,
	chatService *chat.Service,
	knowledgeService *knowledge.Service,
	authService *auth.UserAuthService,
	userService *user.Service,
	jwtManager *auth.JWTManager,
	apiKeyValidator *auth.APIKeyValidator,
) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// NewAPI Integration
	newAPIClient := newapi.NewClient(cfg.NewAPIBaseURL, cfg.NewAPIKey)

	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService)
	billingHandler := handlers.NewBillingHandler(billingService)
	proxyHandler := handlers.NewProxyHandler(billingService, usageService, newAPIClient)
	projectHandler := handlers.NewProjectHandler(projectService, knowledgeService)
	chatHandler := handlers.NewChatSessionHandler(chatService)
	knowledgeHandler := handlers.NewKnowledgeHandler(knowledgeService)
	authHandler := handlers.NewAuthHandler(authService, jwtManager)
	userHandler := handlers.NewUserHandler(userService, authService)
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", middleware.UserAuth(jwtManager), authHandler.Me)

		protected := api.Group("")
		protected.Use(middleware.UserAuth(jwtManager))
		protected.GET("/projects", projectHandler.List)
		protected.POST("/projects", projectHandler.Create)
		protected.GET("/projects/stats", projectHandler.Stats)
		protected.GET("/projects/:id", projectHandler.Get)
		protected.PATCH("/projects/:id", projectHandler.Update)
		protected.DELETE("/projects/:id", projectHandler.Delete)
		protected.GET("/projects/:id/conversations", chatHandler.ListConversations)
		protected.POST("/projects/:id/conversations", chatHandler.CreateConversation)
		protected.GET("/conversations/:conversationId/messages", chatHandler.ListMessages)
		protected.POST("/conversations/:conversationId/messages", chatHandler.CreateMessage)

		protected.GET("/knowledge-bases", knowledgeHandler.ListBases)
		protected.POST("/knowledge-bases", knowledgeHandler.CreateBase)
		protected.GET("/knowledge-bases/:id", knowledgeHandler.GetBase)
		protected.PATCH("/knowledge-bases/:id", knowledgeHandler.UpdateBase)
		protected.DELETE("/knowledge-bases/:id", knowledgeHandler.DeleteBase)
		protected.GET("/knowledge-bases/:id/documents", knowledgeHandler.ListDocuments)
		protected.POST("/knowledge-bases/:id/documents", knowledgeHandler.CreateDocument)
		protected.DELETE("/knowledge-bases/:id/documents/:docId", knowledgeHandler.DeleteDocument)
		protected.GET("/knowledge-bases/:id/documents/:docId/download", knowledgeHandler.DownloadDocument)

		protected.POST("/keys", apiKeyHandler.Create)
		protected.GET("/keys", apiKeyHandler.List)
		protected.POST("/keys/:id/disable", apiKeyHandler.Disable)
		protected.PUT("/keys/:id/scopes", apiKeyHandler.UpdateScopes)
		protected.DELETE("/keys/:id", apiKeyHandler.Delete)

		protected.POST("/billing/hold", billingHandler.Hold)
		protected.POST("/billing/capture", billingHandler.Capture)
		protected.POST("/billing/release", billingHandler.Release)

		protected.GET("/users/me", userHandler.GetMe)
		protected.PATCH("/users/me", userHandler.UpdateMe)
		protected.POST("/users/me/password", userHandler.ChangePassword)
	}

	// Group for AI models (Standard OpenAI format)
	// We forward everything under /v1 to NewAPI
	// This covers /v1/chat/completions, /v1/models, etc.
	v1 := r.Group("/v1")
	{
		v1.Use(middleware.APIKeyAuth(apiKeyValidator))
		// Use Any to match all methods (GET, POST, etc.)
		// /*path will capture the rest of the path
		v1.Any("/*path", proxyHandler.Handle)
	}
}
