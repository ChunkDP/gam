package auth

import (
	"fmt"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var enforcer *casbin.Enforcer

func InitCasbin(db *gorm.DB, rootDir string) (*casbin.Enforcer, error) {
	// 使用 MySQL 适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MySQL adapter: %v", err)
	}

	configPath := filepath.Join(rootDir, "config/rbac_model.conf")

	e, err := casbin.NewEnforcer(configPath, a)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %v", err)
	}

	// 加载策略
	if err := e.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %v", err)
	}

	enforcer = e
	return e, nil
}

func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// 添加权限策略
func AddPolicy(role, path, method string) (bool, error) {
	return enforcer.AddPolicy(role, path, method)
}

// 删除权限策略
func RemovePolicy(role, path, method string) (bool, error) {

	return enforcer.RemovePolicy(role, path, method)
}

// RemoveFilteredPolicy 删除符合条件的权限策略
func RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error) {
	return enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
}

// 添加角色继承关系
func AddRoleForUser(user, role string) (bool, error) {
	return enforcer.AddGroupingPolicy(user, role)
}
