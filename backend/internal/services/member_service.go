package services

import (
	"fmt"
	"normaladmin/backend/internal/models"

	"gorm.io/gorm"
)

type MemberService interface {
	BaseCRUD[models.Member] // 组合基础CRUD接口

	CheckMemberFieldUnique(field, value string, excludeID uint) (bool, error)
}

type memberService struct {
	BaseCRUD[models.Member]
	db *gorm.DB // 保存db实例用于特有方法
}

func NewMemberService(db *gorm.DB, base BaseCRUD[models.Member]) MemberService {

	return &memberService{
		BaseCRUD: base, // 使用装饰后的服务
		db:       db,   // 保存db实例
	}
}

func (s *memberService) CheckMemberFieldUnique(field, value string, excludeID uint) (bool, error) {
	var count int64
	db := s.db.Model(&models.Member{})

	switch field {
	case "username":
		db = db.Where("username = ?", value)
	case "mobile":
		db = db.Where("mobile = ?", value)
	case "email":
		db = db.Where("email = ?", value)
	default:
		return false, fmt.Errorf("invalid field")
	}

	if excludeID != 0 {
		db = db.Where("id != ?", excludeID)
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}
