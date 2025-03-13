package services

import (
	"normaladmin/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService interface {
	BaseCRUD[models.Admin]

	CheckAdminFieldUnique(field, value string, excludeID uint) (bool, error)
	UpdatePassword(id uint, oldPassword, newPassword string) error
}

type adminService struct {
	BaseCRUD[models.Admin]
	db *gorm.DB
}

func NewAdminService(db *gorm.DB, base BaseCRUD[models.Admin]) AdminService {

	return &adminService{
		BaseCRUD: base, // 使用装饰后的服务
		db:       db,   // 保存db实例
	}
}

func (s *adminService) CheckAdminFieldUnique(field, value string, excludeID uint) (bool, error) {
	var count int64
	db := s.db.Model(&models.Admin{}).Where(field+" = ?", value)
	if excludeID > 0 {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *adminService) UpdatePassword(id uint, oldPassword, newPassword string) error {
	var admin models.Admin
	if err := s.db.First(&admin, id).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&admin).Update("password", string(hashedPassword)).Error
}
