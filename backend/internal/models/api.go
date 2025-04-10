package models

import (
	"gorm.io/gorm"
)

// API API管理模型
type API struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Name        string `json:"name" gorm:"not null;comment:API名称"`
	Path        string `json:"path" gorm:"not null;comment:API路径"`
	Method      string `json:"method" gorm:"not null;comment:请求方法"`
	Description string `json:"description" gorm:"comment:API描述"`
	Group       string `json:"group" gorm:"comment:API分组"`
	Status      int    `json:"status" gorm:"default:1;comment:状态(1:启用 0:禁用)"`
	Parameters  string `json:"parameters" gorm:"type:text;comment:请求参数(JSON格式)"`
	Response    string `json:"response" gorm:"type:text;comment:响应示例(JSON格式)"`
	gorm.Model  `swaggerignore:"true"`
}

// TableName 指定表名
func (API) TableName() string {
	return "apis"
}
