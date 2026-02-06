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
		&model.Org{},
		&model.OrgMember{},
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
	)
}

func dropAll(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.KnowledgeDocument{},
		&model.KnowledgeBase{},
		&model.AuditLog{},
		&model.Message{},
		&model.Conversation{},
		&model.UsageRecord{},
		&model.Transaction{},
		&model.Wallet{},
		&model.ProjectDocument{},
		&model.ProjectSkill{},
		&model.ProjectWorkflow{},
		&model.Project{},
		&model.OrgMember{},
		&model.Org{},
		&model.UserSettings{},
		&model.UserProfile{},
		&model.User{},
	)
}
