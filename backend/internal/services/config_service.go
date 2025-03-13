package services

import (
	"fmt"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/cache"
	"normaladmin/backend/pkg/utils/encrypt"
	"time"

	"gorm.io/gorm"
)

type ConfigService interface {
	GetConfigGroups() ([]models.ConfigGroup, error)
	GetConfigItems(groupID int64) ([]models.ConfigItem, error)
	GetConfigValue(groupKey, itemKey string) (string, error)
	UpdateConfigValue(groupID int64, itemKey string, value string) error
	BatchUpdateConfigs(groupID int64, configs map[string]string) error
}

type configService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) ConfigService {
	return &configService{
		db: db,
	}
}

// GetConfigGroups 获取所有配置组
func (s *configService) GetConfigGroups() ([]models.ConfigGroup, error) {
	var groups []models.ConfigGroup
	err := s.db.Where("status = ?", 1).Order("sort_order").Find(&groups).Error
	return groups, err
}

// GetConfigItems 获取配置组的所有配置项
func (s *configService) GetConfigItems(groupID int64) ([]models.ConfigItem, error) {
	var items []models.ConfigItem
	err := s.db.Where("group_id = ?", groupID).Order("sort_order").Find(&items).Error
	if err != nil {
		return nil, err
	}

	// 解密敏感配置
	for i := range items {
		if items[i].Encrypted == 1 {
			items[i].ItemValue = encrypt.Decrypt(items[i].ItemValue)
		}
	}

	// 验证每个配置项的 options 是否为有效的 JSON
	// for i := range items {
	// 	if items[i].Options == "" {
	// 		items[i].Options = "[]"
	// 	}
	// 	// 验证是否为有效的 JSON
	// 	var temp interface{}
	// 	if err := json.Unmarshal([]byte(items[i].Options), &temp); err != nil {
	// 		items[i].Options = "[]" // 如果无效则设置为空数组
	// 	}
	// }

	return items, nil
}

// GetConfigValue 获取配置值
func (s *configService) GetConfigValue(groupKey, itemKey string) (string, error) {
	cacheKey := fmt.Sprintf("config:%s:%s", groupKey, itemKey)

	// 尝试从缓存获取
	if value, err := cache.Get(cacheKey); err == nil {
		return value, nil
	}

	// 从数据库获取
	var item models.ConfigItem
	err := s.db.Joins("JOIN config_groups ON config_items.group_id = config_groups.id").
		Where("config_groups.config_key = ? AND config_items.item_key = ?", groupKey, itemKey).
		First(&item).Error
	if err != nil {
		return "", err
	}

	var value string
	if item.Encrypted == 1 {
		value = encrypt.Decrypt(item.ItemValue)
	} else {
		value = item.ItemValue
	}

	// 设置缓存
	cache.Set(cacheKey, value, time.Hour)

	return value, nil
}

// UpdateConfigValue 更新配置值
func (s *configService) UpdateConfigValue(groupID int64, itemKey string, value string) error {
	var item models.ConfigItem
	err := s.db.Where("group_id = ? AND item_key = ?", groupID, itemKey).First(&item).Error
	if err != nil {
		return err
	}

	// 加密敏感配置
	if item.Encrypted == 1 {
		value = encrypt.Encrypt(value)
	}

	// 更新数据库
	err = s.db.Model(&item).Update("item_value", value).Error
	if err != nil {
		return err
	}

	// 删除缓存
	var group models.ConfigGroup
	s.db.First(&group, groupID)
	cacheKey := fmt.Sprintf("config:%s:%s", group.ConfigKey, itemKey)
	cache.Delete(cacheKey)

	return nil
}

// BatchUpdateConfigs 批量更新配置
func (s *configService) BatchUpdateConfigs(groupID int64, configs map[string]string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for itemKey, value := range configs {
			var item models.ConfigItem
			err := tx.Where("group_id = ? AND item_key = ?", groupID, itemKey).First(&item).Error
			if err != nil {
				return err
			}

			if item.Encrypted == 1 {
				value = encrypt.Encrypt(value)
			}

			if err := tx.Model(&item).Update("item_value", value).Error; err != nil {
				return err
			}
		}

		// 删除组下所有配置的缓存
		var group models.ConfigGroup
		if err := tx.First(&group, groupID).Error; err != nil {
			return err
		}
		cache.Delete(fmt.Sprintf("config:%s:*", group.ConfigKey))

		return nil
	})
}
