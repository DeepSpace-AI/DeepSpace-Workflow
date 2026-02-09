package main

import (
	"log"

	docs "deepspace/cmd/gateway/docs"
	"deepspace/internal/api"
	"deepspace/internal/api/middleware"
	"deepspace/internal/config"
	"deepspace/internal/pkg/db"
	"deepspace/internal/repo"
	"deepspace/internal/service/auth"
	"deepspace/internal/service/billing"
	"deepspace/internal/service/chat"
	"deepspace/internal/service/email"
	"deepspace/internal/service/knowledge"
	modelservice "deepspace/internal/service/model"
	"deepspace/internal/service/passwordreset"
	planservice "deepspace/internal/service/plan"
	"deepspace/internal/service/project"
	"deepspace/internal/service/projectdocument"
	"deepspace/internal/service/projectskill"
	"deepspace/internal/service/projectworkflow"
	"deepspace/internal/service/risk"
	"deepspace/internal/service/usage"
	"deepspace/internal/service/user"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// @title DeepSpace Gateway API
	// @version 1.0
	// @description DeepSpace Gateway 接口文档
	// @BasePath /api
	log.Println("Starting DeepSpace Gateway...")

	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	dbConn, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	if sqlDB, err := dbConn.DB(); err == nil {
		defer sqlDB.Close()
	}
	if err := db.AutoMigrate(dbConn, cfg); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	billingRepo := repo.NewBillingRepo(dbConn)
	billingService := billing.New(dbConn, billingRepo)
	usageRepo := repo.NewUsageRepo(dbConn)
	usageService := usage.New(usageRepo)
	projectRepo := repo.NewProjectRepo(dbConn)
	projectService := project.New(projectRepo)
	projectDocumentRepo := repo.NewProjectDocumentRepo(dbConn)
	projectDocumentService := projectdocument.New(projectDocumentRepo)
	projectSkillRepo := repo.NewProjectSkillRepo(dbConn)
	projectSkillService := projectskill.New(projectSkillRepo)
	projectWorkflowRepo := repo.NewProjectWorkflowRepo(dbConn)
	projectWorkflowService := projectworkflow.New(projectWorkflowRepo)
	chatRepo := repo.NewChatRepo(dbConn)
	chatService := chat.New(chatRepo)
	userRepo := repo.NewUserRepo(dbConn)
	userProfileRepo := repo.NewUserProfileRepo(dbConn)
	userSettingsRepo := repo.NewUserSettingsRepo(dbConn)
	jwtManager := &auth.JWTManager{
		Secret:       []byte(cfg.JWTSecret),
		Issuer:       cfg.JWTIssuer,
		ExpiresIn:    cfg.JWTExpiresIn,
		CookieName:   cfg.JWTCookieName,
		CookieSecure: cfg.JWTCookieSecure,
	}
	userAuthService := auth.NewUserAuthService(userRepo, jwtManager)
	userService := user.New(userRepo, userProfileRepo, userSettingsRepo)
	knowledgeRepo := repo.NewKnowledgeRepo(dbConn)
	knowledgeService := knowledge.New(knowledgeRepo, projectRepo, cfg.KBStoragePath, cfg.KBMaxUploadBytes(), cfg.KBAllowedMIME)
	modelRepo := repo.NewModelRepo(dbConn)
	modelService := modelservice.New(modelRepo)
	planRepo := repo.NewPlanRepo(dbConn)
	planSubscriptionRepo := repo.NewPlanSubscriptionRepo(dbConn)
	planUsageRepo := repo.NewPlanUsageRepo(dbConn)
	planService := planservice.New(planRepo, planSubscriptionRepo, planUsageRepo)
	riskPolicyRepo := repo.NewRiskPolicyRepo(dbConn)
	riskRateRepo := repo.NewRateLimitRepo(dbConn)
	riskIPRepo := repo.NewIPRuleRepo(dbConn)
	riskBudgetRepo := repo.NewBudgetCapRepo(dbConn)
	riskService := risk.New(riskPolicyRepo, riskRateRepo, riskIPRepo, riskBudgetRepo)
	emailService, err := email.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init email service: %v", err)
	}
	passwordResetService, err := passwordreset.New(cfg, userRepo, userProfileRepo, emailService)
	if err != nil {
		log.Fatalf("Failed to init password reset service: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Swagger 基础信息
	docs.SwaggerInfo.Title = "DeepSpace Gateway API"
	docs.SwaggerInfo.Description = "DeepSpace Gateway 接口文档"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"

	r.Use(middleware.TraceID())
	r.Use(middleware.ProjectContext())
	r.Use(middleware.ErrorHandler())
	r.Use(func(c *gin.Context) {
		// Backward-compat: if client mistakenly calls /v1/api/*, strip /v1.
		if strings.HasPrefix(c.Request.URL.Path, "/v1/api/") {
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/v1")
		}
		c.Next()
	})

	// CORS configuration - allow all for now, or restrict as needed
	r.Use(cors.Default())

	// Setup Routes
	api.SetupRoutes(r, cfg, billingService, usageService, projectService, chatService, emailService, knowledgeService, modelService, planService, projectDocumentService, projectSkillService, projectWorkflowService, userAuthService, passwordResetService, userService, riskService, jwtManager)

	log.Printf("Gateway running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
