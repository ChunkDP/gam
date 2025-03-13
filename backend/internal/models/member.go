package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Member 会员模型
type Member struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Username      string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password      string         `gorm:"size:100;not null" json:"-"`
	Nickname      string         `gorm:"size:50" json:"nickname"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	Mobile        string         `gorm:"uniqueIndex;size:20" json:"mobile"`
	Email         string         `gorm:"uniqueIndex;size:100" json:"email"`
	Gender        *int           `gorm:"default:0" json:"gender"` // 0-未知 1-男 2-女
	Birthday      *time.Time     `json:"birthday"`
	LevelID       uint           `gorm:"default:1" json:"level_id"`
	Points        int            `gorm:"default:0" json:"points"`
	Status        *int           `gorm:"default:1" json:"status"` // 0-禁用 1-启用 2-黑名单
	LastLoginTime *time.Time     `json:"last_login_time"`
	LastLoginIP   string         `gorm:"size:50" json:"last_login_ip"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

// MemberLevel 会员等级
type MemberLevel struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Icon        string    `gorm:"size:255" json:"icon"`
	MinPoints   int       `gorm:"not null;default:0" json:"min_points"`
	MaxPoints   int       `gorm:"not null;default:0" json:"max_points"`
	Discount    float64   `gorm:"type:decimal(3,2);default:1.00" json:"discount"`
	Description string    `gorm:"size:255" json:"description"`
	Status      *int      `gorm:"default:1" json:"status"` // 0-禁用 1-启用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MemberPointsLog 积分记录
type MemberPointsLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	MemberID  uint      `gorm:"not null" json:"member_id"`
	Points    int       `gorm:"not null" json:"points"`
	Type      int       `gorm:"not null" json:"type"` // 1-获得 2-消费 3-过期 4-管理员调整
	Source    string    `gorm:"size:50;not null" json:"source"`
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

// MemberTag 会员标签
type MemberTag struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Type      int       `json:"type" gorm:"default:1"`   // 1-系统标签 2-自定义标签
	Status    *int      `json:"status" gorm:"default:1"` // 0-禁用 1-启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type MemberTagRelation struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	MemberID  uint      `gorm:"not null" json:"member_id"`
	TagID     uint      `gorm:"not null" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Member) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// SetPassword 设置密码（会自动加密）
func (a *Member) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
}

// BeforeCreate GORM 的钩子，在创建记录前自动加密密码
func (a *Member) BeforeCreate(tx *gorm.DB) error {
	if a.Password != "" {
		return a.SetPassword(a.Password)
	}
	return nil
}

// BeforeUpdate GORM 的钩子，在更新记录前自动加密密码（如果密码被修改）
func (a *Member) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		return a.SetPassword(a.Password)
	}
	return nil
}
