package models

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID         uint   `json:"id" gorm:"primarykey"`
	ParentID   uint   `json:"parent_id" gorm:"column:parent_id;default:0"`
	Title      string `json:"title" gorm:"not null"`
	Name       string `json:"name" gorm:"not null"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort" gorm:"column:sort;default:0"`
	ParentName string `json:"parent_name" gorm:"column:parent_name;default:''"`
	IsHidden   bool   `json:"is_hidden" gorm:"column:is_hidden;default:false"`
	Type       string `json:"type" gorm:"column:type;default:menu"` // menu or button
	Permission string `json:"permission"`
	ApiMethod  string `json:"api_method"` // 新增字段
	ApiPath    string `json:"api_path"`   // 新增字段
	Status     *int   `json:"status" gorm:"default:1"`
	Children   []Menu `json:"children" gorm:"-"`
	gorm.Model `swaggerignore:"true"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoleID    uint      `json:"role_id" gorm:"column:role_id;not null;comment:角色ID"`
	MenuID    uint      `json:"menu_id" gorm:"column:menu_id;not null;comment:菜单ID"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}
