package migrations

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/models"

	"gorm.io/gorm"
)

// 按顺序注册所有迁移
func init() {
	// 1. 创建基础表
	registerBaseTables()
	// 2. 初始化超级管理员权限
	initSuperAdminPermissions()
}

// registerBaseTables 注册基础表迁移
func registerBaseTables() {
	database.RegisterMigration("001_create_base_tables", func(db *gorm.DB) error {
		return db.AutoMigrate()
	})
}

// registerSystemConfigs 注册系统配置迁移
func initSuperAdminPermissions() {
	database.RegisterMigration("002_init_super_admin_permissions", func(db *gorm.DB) error {
		// 获取所有菜单ID
		var menus []models.Menu
		if err := db.Find(&menus).Error; err != nil {
			return err
		}

		// 为超级管理员(role_id=1)分配所有菜单权限
		var roleMenus []models.RoleMenu
		for _, menu := range menus {
			roleMenus = append(roleMenus, models.RoleMenu{
				RoleID: 1, // 超级管理员角色ID
				MenuID: menu.ID,
			})
		}

		// 批量插入权限记录
		if len(roleMenus) > 0 {
			if err := db.Create(&roleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
