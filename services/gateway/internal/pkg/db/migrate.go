package db

import (
	"deepspace/internal/config"
	"deepspace/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, cfg *config.Config) error {
	if !cfg.DBAutoMigrate {
		return nil
	}

	if cfg.DBResetOnStart {
		if err := dropAll(db); err != nil {
			return err
		}
	}

	return db.AutoMigrate(
		&model.User{},
		&model.UserProfile{},
		&model.UserSettings{},
		&model.Project{},
		&model.ProjectDocument{},
		&model.ProjectSkill{},
		&model.ProjectWorkflow{},
		&model.Wallet{},
		&model.Transaction{},
		&model.UsageRecord{},
		&model.Conversation{},
		&model.Message{},
		&model.AuditLog{},
		&model.KnowledgeBase{},
		&model.KnowledgeDocument{},
		&model.Model{},
		&model.Plan{},
		&model.PlanModelPrice{},
		&model.PlanSubscription{},
	)
}

func dropAll(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.User{},
		&model.UserProfile{},
		&model.UserSettings{},
		&model.Project{},
		&model.ProjectDocument{},
		&model.ProjectSkill{},
		&model.ProjectWorkflow{},
		&model.Wallet{},
		&model.Transaction{},
		&model.UsageRecord{},
		&model.Conversation{},
		&model.Message{},
		&model.AuditLog{},
		&model.KnowledgeBase{},
		&model.KnowledgeDocument{},
		&model.Model{},
		&model.PlanSubscription{},
		&model.PlanModelPrice{},
		&model.Plan{},
	)
}
