package database

import (
	"fmt"

	"kira-go/internal/config"
	"kira-go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect 初始化数据库连接
func Connect() error {
	var err error
	DB, err = gorm.Open(postgres.Open(config.C.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 生产用 logger.Error
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 SQL 数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not connected, call Connect() first")
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.GitHubUser{},
		&models.Post{},
		&models.Category{},
		&models.Tag{},
		&models.PostTag{},
		&models.Album{},
		&models.Photo{},
		&models.Project{},
		&models.SiteConfig{},
		&models.Visitor{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}
