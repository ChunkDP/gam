package services

import (
	"fmt"
	"normaladmin/backend/internal/models"

	"gorm.io/gorm"
)

type MenuService interface {
	GetMenuTree() ([]models.Menu, error)
	CreateMenu(menu *models.Menu) error
	UpdateMenu(id uint, menu *models.Menu) error
	DeleteMenu(id uint) error
	GetMenus(query map[string]interface{}) ([]models.Menu, error)
	Update(id uint, data interface{}) error
}

type menuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) MenuService {
	return &menuService{db: db}
}

func (s *menuService) GetMenuTree() ([]models.Menu, error) {
	var menus []models.Menu
	if err := s.db.Order("sort").Find(&menus).Error; err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

func (s *menuService) CreateMenu(menu *models.Menu) error {
	if menu.ParentID != 0 {
		var parent models.Menu
		if err := s.db.First(&parent, menu.ParentID).Error; err != nil {
			return fmt.Errorf("parent menu not found")
		}
		menu.ParentName = parent.Name
	}
	return s.db.Create(menu).Error
}

func (s *menuService) UpdateMenu(id uint, menu *models.Menu) error {
	var existing models.Menu
	if err := s.db.First(&existing, id).Error; err != nil {
		return err
	}

	if menu.ParentID != existing.ParentID {
		if err := s.validateMenuHierarchy(id, menu.ParentID); err != nil {
			return err
		}
	}

	return s.db.Save(menu).Error
}

func (s *menuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var childCount int64
	if err := s.db.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete menu with child items")
	}

	// 检查是否被角色使用
	var roleMenuCount int64
	if err := s.db.Model(&models.RoleMenu{}).Where("menu_id = ?", id).Count(&roleMenuCount).Error; err != nil {
		return err
	}
	if roleMenuCount > 0 {
		return fmt.Errorf("cannot delete menu that is assigned to roles")
	}

	return s.db.Delete(&models.Menu{}, id).Error
}

func (s *menuService) GetMenus(query map[string]interface{}) ([]models.Menu, error) {
	var menus []models.Menu
	db := s.db.Order("sort")

	for key, value := range query {
		if str, ok := value.(string); ok && str != "" {
			db = db.Where(key+" LIKE ?", "%"+str+"%")
		}
	}

	if err := db.Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
func (s *menuService) Update(id uint, data interface{}) error {
	var model models.Menu
	if err := s.db.First(&model, id).Error; err != nil {
		return err
	}

	return s.db.Model(&models.Menu{}).Where("id = ?", id).Updates(data).Error
}

// 辅助方法
func (s *menuService) validateMenuHierarchy(menuID uint, parentID uint) error {
	if parentID == 0 {
		return nil
	}

	if menuID == parentID {
		return fmt.Errorf("cannot select self as parent")
	}

	var parent models.Menu
	if err := s.db.First(&parent, parentID).Error; err != nil {
		return fmt.Errorf("parent menu not found")
	}

	isChild, err := s.isChildMenu(menuID, parentID)
	if err != nil {
		return err
	}
	if isChild {
		return fmt.Errorf("cannot select child menu as parent")
	}

	return nil
}

func (s *menuService) isChildMenu(menuID uint, targetID uint) (bool, error) {
	var children []models.Menu
	if err := s.db.Where("parent_id = ?", menuID).Find(&children).Error; err != nil {
		return false, err
	}

	for _, child := range children {
		if child.ID == targetID {
			return true, nil
		}
		isChild, err := s.isChildMenu(child.ID, targetID)
		if err != nil {
			return false, err
		}
		if isChild {
			return true, nil
		}
	}

	return false, nil
}

func buildMenuTree(menus []models.Menu, parentID uint) []models.Menu {
	var tree []models.Menu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = buildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}
