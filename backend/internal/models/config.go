package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// ConfigGroup 配置组
type ConfigGroup struct {
	gorm.Model     `swaggerignore:"true"`
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"` // 将 ID 字段的 JSON 标签设置为 "id"
	ConfigKey      string `gorm:"column:config_key;type:varchar(50);uniqueIndex;not null" json:"config_key"`
	ConfigName     string `gorm:"column:config_name;type:varchar(100);not null" json:"config_name"`
	Description    string `gorm:"column:description;type:text" json:"description"`
	Icon           string `gorm:"column:icon;type:varchar(50)" json:"icon"`
	SortOrder      int    `gorm:"column:sort_order;default:0" json:"sort_order"`
	Status         int8   `gorm:"column:status;default:1" json:"status"`
	ParentID       int64  `gorm:"column:parent_id;default:0" json:"parent_id"`
	PermissionCode string `gorm:"column:permission_code;type:varchar(100)" json:"permission_code"`
	AdminRoles     string `gorm:"column:admin_roles;type:varchar(255)" json:"admin_roles"`
	ViewRoles      string `gorm:"column:view_roles;type:varchar(255)" json:"view_roles"`
	IsSystem       int8   `gorm:"column:is_system;default:0" json:"is_system"`
}

// ConfigItem 配置项
type ConfigItem struct {
	gorm.Model       `swaggerignore:"true"`
	ID               uint        `json:"id" gorm:"primaryKey;autoIncrement"` // 将 ID 字段的 JSON 标签设置为 "id"
	GroupID          int64       `gorm:"column:group_id;not null" json:"group_id"`
	ItemKey          string      `gorm:"column:item_key;type:varchar(50);not null" json:"item_key"`
	ItemName         string      `gorm:"column:item_name;type:varchar(100);not null" json:"item_name"`
	ItemValue        string      `gorm:"column:item_value;type:text" json:"item_value"`
	ValueType        string      `gorm:"column:value_type;type:varchar(20);not null" json:"value_type"`
	Description      string      `gorm:"column:description;type:text" json:"description"`
	SortOrder        int         `gorm:"column:sort_order;default:0" json:"sort_order"`
	Required         int8        `gorm:"column:required;default:0" json:"required"`
	Options          OptionsJSON `gorm:"column:options;type:text" json:"options"`
	DependsOn        string      `gorm:"column:depends_on;type:varchar(255)" json:"depends_on"`
	VisibleCondition string      `gorm:"column:visible_condition;type:varchar(255)" json:"visible_condition"`
	Encrypted        int8        `gorm:"column:encrypted;default:0" json:"encrypted"`
}

// OptionsJSON 自定义JSON类型
type OptionsJSON []Option

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Scan 实现 sql.Scanner 接口
func (o *OptionsJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, o)
}

// Value 实现 driver.Valuer 接口
func (o OptionsJSON) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return json.Marshal(o)
}
