package main

import (
	"log"
	"strings"

	"deepspace/internal/api"
	"deepspace/internal/api/middleware"
	"deepspace/internal/config"
	"deepspace/internal/pkg/db"
	"deepspace/internal/repo"
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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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
	orgRepo := repo.NewOrgRepo(dbConn)
	userProfileRepo := repo.NewUserProfileRepo(dbConn)
	userSettingsRepo := repo.NewUserSettingsRepo(dbConn)
	jwtManager := &auth.JWTManager{
		Secret:       []byte(cfg.JWTSecret),
		Issuer:       cfg.JWTIssuer,
		ExpiresIn:    cfg.JWTExpiresIn,
		CookieName:   cfg.JWTCookieName,
		CookieSecure: cfg.JWTCookieSecure,
	}
	userAuthService := auth.NewUserAuthService(userRepo, orgRepo, jwtManager)
	userService := user.New(userRepo, userProfileRepo, userSettingsRepo)
	knowledgeRepo := repo.NewKnowledgeRepo(dbConn)
	knowledgeService := knowledge.New(knowledgeRepo, projectRepo, cfg.KBStoragePath, cfg.KBMaxUploadBytes(), cfg.KBAllowedMIME)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(middleware.TraceID())
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
	api.SetupRoutes(r, cfg, billingService, usageService, projectService, chatService, knowledgeService, projectDocumentService, projectSkillService, projectWorkflowService, userAuthService, userService, jwtManager)

	log.Printf("Gateway running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
