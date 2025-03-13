package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model  `swaggerignore:"true"`
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"` // 将 ID 字段的 JSON 标签设置为 "id"
	Name        string `json:"name" gorm:"type:varchar(50);unique;not null"`
	Code        string `json:"code" gorm:"type:varchar(50);unique;not null"`
	Description string `json:"description" gorm:"type:varchar(200)"`
	IsPreset    bool   `json:"is_preset" gorm:"column:is_preset;default:false"`                                                                                // 是否为系统预设角色
	Status      *int   `json:"status" gorm:"type:tinyint;default:1;comment:'状态 1:启用 2:禁用'"`                                                                    // 状态
	Sort        int    `json:"sort" gorm:"type:int;default:0;comment:'排序（值越小越靠前）'"`                                                                            // 排序
	Remark      string `json:"remark" gorm:"type:varchar(500);comment:'备注信息'"`                                                                                 // 备注
	DataScope   int    `json:"data_scope" gorm:"column:data_scope;type:tinyint;default:1;comment:'数据范围（1：全部数据权限 2：自定义数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）'"` // 数据权限范围
}
