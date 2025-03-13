package models

import (
	"time"
)

// NotificationType 通知类型
type NotificationType struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"size:50;not null;comment:类型名称"`
	Code        string    `json:"code" gorm:"size:50;not null;uniqueIndex;comment:类型编码"`
	Description string    `json:"description" gorm:"size:255;comment:类型描述"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Notification 通知
type Notification struct {
	ID             uint             `json:"id" gorm:"primarykey"`
	Title          string           `json:"title" gorm:"size:100;not null;comment:通知标题"`
	Content        string           `json:"content" gorm:"type:text;comment:通知内容"`
	TypeID         uint             `json:"type_id" gorm:"not null;comment:通知类型ID"`
	Type           NotificationType `json:"type" gorm:"foreignKey:TypeID"`
	Level          int              `json:"level" gorm:"default:1;comment:重要程度(1普通 2重要 3紧急)"`
	Status         int              `json:"status" gorm:"default:0;comment:状态(0草稿 1已发布 2已撤回)"`
	SenderID       uint             `json:"sender_id" gorm:"comment:发送者ID"`
	SenderName     string           `json:"sender_name" gorm:"size:50;comment:发送者名称"`
	PublishTime    *time.Time       `json:"publish_time" gorm:"comment:发布时间"`
	ExpirationTime string           `json:"expiration_time" gorm:"comment:过期时间"`
	ReceiverType   string           `json:"receiver_type" gorm:"size:20;not null;default:'all';comment:接收者类型(all/members/admins)"`
	ReadCount      int              `json:"read_count" gorm:"default:0;comment:已读数量"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// NotificationReceiver 通知接收记录
type NotificationReceiver struct {
	ID             uint         `json:"id" gorm:"primarykey"`
	NotificationID uint         `json:"notification_id" gorm:"not null;comment:通知ID"`
	UserID         uint         `json:"user_id" gorm:"not null;comment:用户ID"`
	UserType       string       `json:"user_type" gorm:"type:varchar(20);not null;comment:用户类型(admin/member)"`
	UserName       string       `json:"user_name" gorm:"size:50;comment:用户名称"`
	IsRead         bool         `json:"is_read" gorm:"default:false;comment:是否已读"`
	ReadTime       *time.Time   `json:"read_time" gorm:"comment:阅读时间"`
	IsDeleted      bool         `json:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Notification   Notification `json:"notification" gorm:"foreignKey:NotificationID"`
}

// NotificationQuery 通知查询参数
type NotificationQuery struct {
	Title        string `form:"title"`
	TypeID       uint   `form:"type_id"`
	Level        int    `form:"level"`
	Status       int    `form:"status"`
	ReceiverType string `form:"receiver_type"`
	StartTime    string `form:"start_time"`
	EndTime      string `form:"end_time"`
	Page         int    `form:"page" binding:"required,min=1"`
	PageSize     int    `form:"page_size" binding:"required,min=5,max=100"`
}

// UserNotificationQuery 用户通知查询参数
type UserNotificationQuery struct {
	IsRead    *bool  `form:"is_read"`
	TypeID    uint   `form:"type_id"`
	Level     int    `form:"level"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Page      int    `form:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=5,max=100"`
}

// ReceiverQuery 接收者查询参数
type ReceiverQuery struct {
	UserType string `form:"user_type"`
	UserName string `form:"user_name"`
	IsRead   *bool  `form:"is_read"`
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=100"`
}
