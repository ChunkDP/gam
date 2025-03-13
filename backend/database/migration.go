package database

import (
	"fmt"
	"time"

	"normaladmin/backend/pkg/logger"

	"gorm.io/gorm"
)

var migrations []Migration

// Migration 迁移定义
type Migration struct {
	Name string
	Fn   MigrationFunc
}

type MigrationFunc func(*gorm.DB) error

// RegisterMigration 注册迁移函数
func RegisterMigration(name string, fn MigrationFunc) {
	migrations = append(migrations, Migration{
		Name: name,
		Fn:   fn,
	})
}

// Migrate 执行所有迁移
func Migrate() error {
	db := GetDB()

	logger.Info("开始执行迁移", logger.Field("total", len(migrations)))

	// 创建迁移记录表
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// 执行所有未执行的迁移
	for _, migration := range migrations {
		logger.Info("执行迁移",
			logger.Field("name", migration.Name),
		)

		// 检查是否已执行
		var record MigrationRecord
		if err := db.Where("name = ?", migration.Name).First(&record).Error; err == nil {
			// 迁移已执行，跳过
			logger.Info("迁移完成",
				logger.Field("name", migration.Name),
			)
			continue
		}

		// 执行迁移
		if err := migration.Fn(db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Name, err)
		}

		// 记录迁移
		record = MigrationRecord{
			Name:      migration.Name,
			CreatedAt: time.Now().Unix(),
		}
		if err := db.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		logger.Info("迁移完成",
			logger.Field("name", migration.Name),
		)
	}

	return nil
}

// MigrationRecord 迁移记录
type MigrationRecord struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"size:255;not null;uniqueIndex"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}
