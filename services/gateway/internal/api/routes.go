package api

import (
	"net/http"

	"deepspace/internal/api/handlers"
	"deepspace/internal/api/middleware"
	"deepspace/internal/config"
	"deepspace/internal/integrations/newapi"
	"deepspace/internal/service/auth"
	"deepspace/internal/service/billing"
	"deepspace/internal/service/chat"
	"deepspace/internal/service/email"
	"deepspace/internal/service/knowledge"
	modelservice "deepspace/internal/service/model"
	planservice "deepspace/internal/service/plan"
	"deepspace/internal/service/project"
	"deepspace/internal/service/projectdocument"
	"deepspace/internal/service/projectskill"
	"deepspace/internal/service/projectworkflow"
	"deepspace/internal/service/usage"
	"deepspace/internal/service/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	billingService *billing.Service,
	usageService *usage.Service,
	projectService *project.Service,
	chatService *chat.Service,
	emailService *email.Service,
	knowledgeService *knowledge.Service,
	modelService *modelservice.Service,
	planService *planservice.Service,
	projectDocumentService *projectdocument.Service,
	projectSkillService *projectskill.Service,
	projectWorkflowService *projectworkflow.Service,
	authService *auth.UserAuthService,
	userService *user.Service,
	jwtManager *auth.JWTManager,
) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// NewAPI Integration
	newAPIClient := newapi.NewClient(cfg.NewAPIBaseURL, cfg.NewAPIKey)

	billingHandler := handlers.NewBillingHandler(billingService)
	billingViewHandler := handlers.NewBillingViewHandler(billingService, usageService)
	proxyHandler := handlers.NewProxyHandler(billingService, usageService, newAPIClient, modelService, planService)
	projectHandler := handlers.NewProjectHandler(projectService, knowledgeService)
	chatHandler := handlers.NewChatSessionHandler(chatService)
	emailHandler := handlers.NewEmailHandler(emailService)
	knowledgeHandler := handlers.NewKnowledgeHandler(knowledgeService)
	modelHandler := handlers.NewModelHandler(modelService, newAPIClient)
	planHandler := handlers.NewPlanHandler(planService)
	planSubscriptionHandler := handlers.NewPlanSubscriptionHandler(planService)
	projectDocumentHandler := handlers.NewProjectDocumentHandler(projectDocumentService)
	projectSkillHandler := handlers.NewProjectSkillHandler(projectSkillService)
	projectWorkflowHandler := handlers.NewProjectWorkflowHandler(projectWorkflowService)
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
		protected.GET("/projects/:id/documents", projectDocumentHandler.List)
		protected.POST("/projects/:id/documents", projectDocumentHandler.Create)
		protected.GET("/projects/:id/documents/:docId", projectDocumentHandler.Get)
		protected.PATCH("/projects/:id/documents/:docId", projectDocumentHandler.Update)
		protected.DELETE("/projects/:id/documents/:docId", projectDocumentHandler.Delete)
		protected.GET("/projects/:id/skills", projectSkillHandler.List)
		protected.POST("/projects/:id/skills", projectSkillHandler.Create)
		protected.PATCH("/projects/:id/skills/:skillId", projectSkillHandler.Update)
		protected.DELETE("/projects/:id/skills/:skillId", projectSkillHandler.Delete)
		protected.GET("/projects/:id/workflows", projectWorkflowHandler.List)
		protected.POST("/projects/:id/workflows", projectWorkflowHandler.Create)
		protected.PATCH("/projects/:id/workflows/:workflowId", projectWorkflowHandler.Update)
		protected.DELETE("/projects/:id/workflows/:workflowId", projectWorkflowHandler.Delete)
		protected.GET("/projects/:id/conversations", chatHandler.ListConversations)
		protected.POST("/projects/:id/conversations", chatHandler.CreateConversation)
		protected.GET("/conversations", chatHandler.ListStandaloneConversations)
		protected.POST("/conversations", chatHandler.CreateStandaloneConversation)
		protected.GET("/conversations/:conversationId/messages", chatHandler.ListMessages)
		protected.POST("/conversations/:conversationId/messages", chatHandler.CreateMessage)
		protected.PATCH("/conversations/:conversationId", chatHandler.UpdateConversation)
		protected.DELETE("/conversations/:conversationId", chatHandler.DeleteConversation)

		protected.GET("/knowledge-bases", knowledgeHandler.ListBases)
		protected.POST("/knowledge-bases", knowledgeHandler.CreateBase)
		protected.GET("/knowledge-bases/:id", knowledgeHandler.GetBase)
		protected.PATCH("/knowledge-bases/:id", knowledgeHandler.UpdateBase)
		protected.DELETE("/knowledge-bases/:id", knowledgeHandler.DeleteBase)
		protected.GET("/knowledge-bases/:id/documents", knowledgeHandler.ListDocuments)
		protected.POST("/knowledge-bases/:id/documents", knowledgeHandler.CreateDocument)
		protected.DELETE("/knowledge-bases/:id/documents/:docId", knowledgeHandler.DeleteDocument)
		protected.GET("/knowledge-bases/:id/documents/:docId/download", knowledgeHandler.DownloadDocument)

		protected.POST("/billing/hold", billingHandler.Hold)
		protected.POST("/billing/capture", billingHandler.Capture)
		protected.POST("/billing/release", billingHandler.Release)
		protected.GET("/billing/wallet", billingViewHandler.Wallet)
		protected.GET("/billing/usage", billingViewHandler.Usage)

		protected.GET("/users/me", userHandler.GetMe)
		protected.PATCH("/users/me", userHandler.UpdateMe)
		protected.POST("/users/me/password", userHandler.ChangePassword)
		protected.POST("/email/send", emailHandler.Send)
		protected.POST("/email/enqueue", emailHandler.Enqueue)
		protected.GET("/models", modelHandler.List)
		protected.GET("/models/providers", modelHandler.ListProviders)

		// Admin User Management
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireAdmin(userService))
		{
			admin.GET("/users", userHandler.List)
			admin.POST("/users", userHandler.Create)
			admin.GET("/users/:id", userHandler.Get)
			admin.PATCH("/users/:id", userHandler.Update)
			admin.DELETE("/users/:id", userHandler.Delete)

			admin.POST("/models/sync", modelHandler.Sync)
			admin.POST("/models/confirm", modelHandler.ConfirmBatch)
			admin.GET("/models", modelHandler.ListAll)
			admin.GET("/models/providers", modelHandler.ListAllProviders)
			admin.POST("/models/pricing", modelHandler.BatchPricing)
			admin.POST("/models", modelHandler.Create)
			admin.PATCH("/models/:id", modelHandler.Update)
			admin.GET("/plans", planHandler.List)
			admin.POST("/plans", planHandler.Create)
			admin.PATCH("/plans/:id", planHandler.Update)
			admin.POST("/subscriptions", planSubscriptionHandler.Create)
			admin.PATCH("/subscriptions/:id", planSubscriptionHandler.Update)
			admin.GET("/users/:id/subscription", planSubscriptionHandler.GetOrgActive)
		}
	}

	// Group for AI models (Standard OpenAI format)
	// We forward everything under /v1 to NewAPI
	// This covers /v1/chat/completions, /v1/models, etc.
	v1 := r.Group("/v1")
	{
		v1.Use(middleware.UserAuth(jwtManager))
		// Use Any to match all methods (GET, POST, etc.)
		// /*path will capture the rest of the path
		v1.Any("/*path", proxyHandler.Handle)
	}
}
