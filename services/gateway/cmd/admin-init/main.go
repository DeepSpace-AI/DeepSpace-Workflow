package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"deepspace/internal/config"
	"deepspace/internal/model"
	"deepspace/internal/pkg/db"
	"deepspace/internal/repo"
	"deepspace/internal/service/user"

	"gorm.io/gorm"
)

func main() {
	log.Println("开始初始化管理员用户...")

	email := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))
	password := os.Getenv("ADMIN_PASSWORD")
	displayName := strings.TrimSpace(os.Getenv("ADMIN_DISPLAY_NAME"))

	if email == "" || strings.TrimSpace(password) == "" {
		log.Fatal("未提供 ADMIN_EMAIL 或 ADMIN_PASSWORD")
	}

	cfg := config.Load()
	dbConn, err := db.New(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err := db.AutoMigrate(dbConn, cfg); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	ctx := context.Background()
	var existing model.User
	err = dbConn.WithContext(ctx).Where("role = ?", "admin").First(&existing).Error
	if err == nil {
		log.Println("已存在管理员用户，跳过初始化")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("查询管理员失败: %v", err)
	}

	userRepo := repo.NewUserRepo(dbConn)
	profileRepo := repo.NewUserProfileRepo(dbConn)
	settingsRepo := repo.NewUserSettingsRepo(dbConn)
	userService := user.New(userRepo, profileRepo, settingsRepo)

	var profile *user.UpdateProfile
	if displayName != "" {
		name := displayName
		profile = &user.UpdateProfile{DisplayName: &name}
	}

	_, err = userService.Create(ctx, user.CreateInput{
		Email:    email,
		Password: password,
		Role:     "admin",
		Status:   "active",
		Profile:  profile,
	})
	if err != nil {
		if err == user.ErrEmailTaken {
			log.Fatal("管理员邮箱已被占用")
		}
		log.Fatalf("创建管理员失败: %v", err)
	}

	log.Println("管理员用户初始化完成")
}
