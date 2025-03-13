package services

import (
	"fmt"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/auth"

	"gorm.io/gorm"
)

type RoleService interface {
	BaseCRUD[models.Role] // 组合基础CRUD接口
	CheckRoleFieldUnique(field, value string, excludeID uint) (bool, error)
	UpdateRoleMenus(roleID uint, menuIDs []uint) error
	GetRoleMenus(roleID uint) ([]map[string]interface{}, []uint, error)
}

type roleService struct {
	BaseCRUD[models.Role]
	db *gorm.DB
}

func NewRoleService(db *gorm.DB, base BaseCRUD[models.Role]) RoleService {
	return &roleService{
		BaseCRUD: base,
		db:       db,
	}
}

func (s *roleService) CheckRoleFieldUnique(field, value string, excludeID uint) (bool, error) {
	var count int64
	db := s.db.Model(&models.Role{}).Where(field+" = ?", value)
	if excludeID > 0 {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// UpdateRoleMenus 更新角色菜单权限(包含事务处理)
func (s *roleService) UpdateRoleMenus(roleID uint, menuIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除原有的角色-菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}

		// 创建新的角色-菜单关联
		for _, menuID := range menuIDs {
			roleMenu := models.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			}
			if err := tx.Create(&roleMenu).Error; err != nil {
				return err
			}

			// 如果是按钮类型的菜单，同时更新 Casbin 策略
			var menu models.Menu
			if err := tx.First(&menu, menuID).Error; err != nil {
				return err
			}

			if menu.Type == "button" && menu.ApiMethod != "" && menu.ApiPath != "" {
				if _, err := auth.AddPolicy(fmt.Sprint(roleID), menu.ApiPath, menu.ApiMethod); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// GetRoleMenus 获取角色的菜单权限
func (s *roleService) GetRoleMenus(roleID uint) ([]map[string]interface{}, []uint, error) {
	checkedMenus := make([]uint, 0)
	// 获取所有菜单
	var menus []models.Menu
	if err := s.db.Order("sort").Find(&menus).Error; err != nil {
		return nil, nil, err
	}

	// 获取角色已有的菜单权限
	var roleMenus []models.RoleMenu
	if err := s.db.Where("role_id = ?", roleID).Find(&roleMenus).Error; err != nil {
		return nil, nil, err
	}

	for _, rm := range roleMenus {
		checkedMenus = append(checkedMenus, rm.MenuID)
	}
	// 构建带选中状态的菜单树
	menuTree := buildMenuTreeWithChecked(menus, roleMenus)
	return menuTree, checkedMenus, nil
}

// buildMenuTreeWithChecked 构建带选中状态的菜单树
func buildMenuTreeWithChecked(menus []models.Menu, roleMenus []models.RoleMenu) []map[string]interface{} {
	// 创建已选中菜单ID的map，方便查找
	checkedMenus := make(map[uint]bool)
	for _, rm := range roleMenus {
		checkedMenus[rm.MenuID] = true
	}

	return buildMenuTreeWithCheckedHelper(menus, 0, checkedMenus)
}

func buildMenuTreeWithCheckedHelper(menus []models.Menu, parentID uint, checkedMenus map[uint]bool) []map[string]interface{} {
	var tree []map[string]interface{}

	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := map[string]interface{}{
				"id":       menu.ID,
				"title":    menu.Title,
				"checked":  checkedMenus[menu.ID],
				"children": buildMenuTreeWithCheckedHelper(menus, menu.ID, checkedMenus),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
