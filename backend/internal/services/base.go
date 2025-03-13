package services

import (
	"normaladmin/backend/internal/models"

	"gorm.io/gorm"
)

// CRUDService 基础CRUD服务接口
type BaseCRUD[T any] interface {
	GetByID(id uint, opts ...models.QueryOption) (*T, error)
	List(query map[string]interface{}, page, pageSize string, opts ...models.QueryOption) ([]T, int64, error)
	Create(entity *T) error
	Update(id uint, data interface{}) error
	Delete(id uint, hardDelete bool) error
	BatchDelete(ids []uint, hardDelete bool) error
}

// BaseCRUDService 基础实现
type BaseCRUDService[T any] struct {
	db *gorm.DB
}

func NewBaseCRUDService[T any](db *gorm.DB) BaseCRUD[T] {
	return &BaseCRUDService[T]{
		db: db,
	}
}

// GetByID 实现带查询选项的获取方法
func (s *BaseCRUDService[T]) GetByID(id uint, opts ...models.QueryOption) (*T, error) {
	var data T
	db := s.db

	// 应用所有查询选项
	for _, opt := range opts {
		db = opt(db)
	}

	if err := db.First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// List 实现带查询选项的列表方法
func (s *BaseCRUDService[T]) List(query map[string]interface{}, page, pageSize string, opts ...models.QueryOption) ([]T, int64, error) {
	var data []T
	var total int64
	db := s.db

	// 应用所有查询选项
	for _, opt := range opts {
		db = opt(db)
	}

	// 处理查询条件
	for field, value := range query {
		if strValue, ok := value.(string); ok && strValue != "" {
			db = db.Where(field+" LIKE ?", "%"+strValue+"%")
		} else if value != nil {
			db = db.Where(field+" = ?", value)
		}
	}

	// 获取总数
	if err := db.Model(&data).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 使用 Paginate 进行分页查询
	if err := db.Scopes(models.Paginate(page, pageSize)).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *BaseCRUDService[T]) Create(entity *T) error {
	return s.db.Create(entity).Error
}

// Update 支持实体对象或map更新
func (s *BaseCRUDService[T]) Update(id uint, data interface{}) error {
	var model T
	if err := s.db.First(&model, id).Error; err != nil {
		return err
	}

	return s.db.Model(&model).Updates(data).Error
}

// Delete 支持软删除和硬删除
func (s *BaseCRUDService[T]) Delete(id uint, hardDelete bool) error {
	var model T
	if err := s.db.First(&model, id).Error; err != nil {
		return err
	}

	if hardDelete {
		// 硬删除
		return s.db.Unscoped().Delete(&model).Error
	}
	// 软删除
	return s.db.Delete(&model).Error
}

// BatchDelete 批量删除支持
func (s *BaseCRUDService[T]) BatchDelete(ids []uint, hardDelete bool) error {
	var model T
	if hardDelete {
		return s.db.Unscoped().Delete(&model, ids).Error
	}
	return s.db.Delete(&model, ids).Error
}
