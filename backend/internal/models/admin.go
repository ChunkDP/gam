package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Username      string         `json:"username" gorm:"unique;not null"`
	Password      string         `json:"-" gorm:"not null"` // json:"-" 表示不返回密码
	Email         string         `json:"email" gorm:"size:100"`
	Phone         string         `json:"phone" gorm:"size:20"`
	RealName      string         `json:"real_name" gorm:"size:50"`
	Avatar        string         `json:"avatar" gorm:"size:255"`
	RoleID        uint           `json:"role_id"`
	Role          Role           `json:"role" gorm:"foreignKey:RoleID"`
	Status        *int           `json:"status" gorm:"default:1"` // 1-启用 0-禁用
	LastLoginTime *time.Time     `json:"last_login_time"`         // 改为指针类型，允许为 null
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index" swaggerignore:"true"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

// CheckPassword 验证密码是否正确
func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// SetPassword 设置密码（会自动加密）
func (a *Admin) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
}

// BeforeCreate GORM 的钩子，在创建记录前自动加密密码
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.Password != "" {
		return a.SetPassword(a.Password)
	}
	return nil
}

// BeforeUpdate GORM 的钩子，在更新记录前自动加密密码（如果密码被修改）
func (a *Admin) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		return a.SetPassword(a.Password)
	}
	return nil
}
