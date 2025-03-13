package sysconfig

import (
	"encoding/json"
	"fmt"
	"normaladmin/backend/pkg/utils/encrypt"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	instance *SysConfig
	once     sync.Once
)

type ConfigItem struct {
	ItemKey   string `gorm:"column:item_key"`
	ItemValue string `gorm:"column:item_value"`
	ValueType string `gorm:"column:value_type"`
	Encrypted bool   `gorm:"column:encrypted"`
}

type SysConfig struct {
	db       *gorm.DB
	cache    map[string]string
	mutex    sync.RWMutex
	lastLoad time.Time
	cacheTTL time.Duration
}

// loadConfig 从数据库加载配置
func (s *SysConfig) loadConfig() error {
	var items []ConfigItem
	if err := s.db.Table("config_items").
		Select("item_key", "item_value", "value_type", "encrypted").
		Find(&items).Error; err != nil {
		return err
	}

	newCache := make(map[string]string)
	for _, item := range items {
		value := item.ItemValue
		// 如果是加密的配置项，需要解密
		if item.Encrypted {
			decrypted := encrypt.Decrypt(value)

			value = decrypted
		}
		newCache[item.ItemKey] = value
	}

	s.mutex.Lock()
	s.cache = newCache
	s.lastLoad = time.Now()
	s.mutex.Unlock()

	return nil
}

// checkAndReload 检查是否需要重新加载配置
func (s *SysConfig) checkAndReload() {
	s.mutex.RLock()
	expired := time.Since(s.lastLoad) > s.cacheTTL
	s.mutex.RUnlock()

	if expired {
		s.loadConfig()
	}
}

// Get 获取配置值
func (s *SysConfig) Get(key string, defaultValue string) string {
	s.checkAndReload()

	s.mutex.RLock()
	value, exists := s.cache[key]
	s.mutex.RUnlock()

	if !exists {
		return defaultValue
	}
	return value
}

// GetInt 获取整数配置
func (s *SysConfig) GetInt(key string, defaultValue int) int {
	value := s.Get(key, "")
	if value == "" {
		return defaultValue
	}

	intValue := 0
	if err := json.Unmarshal([]byte(value), &intValue); err != nil {
		return defaultValue
	}
	return intValue
}

// GetInt64 获取int64类型配置
func (s *SysConfig) GetInt64(key string, defaultValue int64) int64 {
	value := s.Get(key, fmt.Sprintf("%d", defaultValue))
	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int64Value
}

// GetBool 获取布尔配置
func (s *SysConfig) GetBool(key string, defaultValue bool) bool {
	value := s.Get(key, "")
	if value == "" {
		return defaultValue
	}

	boolValue := false
	if err := json.Unmarshal([]byte(value), &boolValue); err != nil {
		return defaultValue
	}
	return boolValue
}

// GetFloat 获取浮点数配置
func (s *SysConfig) GetFloat(key string, defaultValue float64) float64 {
	value := s.Get(key, "")
	if value == "" {
		return defaultValue
	}

	floatValue := 0.0
	if err := json.Unmarshal([]byte(value), &floatValue); err != nil {
		return defaultValue
	}
	return floatValue
}

// GetArray 获取数组配置
func (s *SysConfig) GetArray(key string, defaultValue []string) []string {
	value := s.Get(key, "")
	if value == "" {
		return defaultValue
	}

	var arrayValue []string
	if err := json.Unmarshal([]byte(value), &arrayValue); err != nil {
		return defaultValue
	}
	return arrayValue
}

// Set 设置配置项
func (s *SysConfig) Set(key string, value interface{}) error {
	// 将值转换为JSON字符串
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 获取配置项信息
	var item ConfigItem
	if err := s.db.Table("config_items").
		Select("encrypted").
		Where("item_key = ?", key).
		First(&item).Error; err != nil {
		return err
	}

	// 如果是加密配置项，需要加密
	valueStr := string(jsonValue)
	if item.Encrypted {
		encrypted := encrypt.Encrypt(valueStr)

		valueStr = encrypted
	}

	// 更新数据库
	if err := s.db.Table("config_items").
		Where("item_key = ?", key).
		Update("item_value", valueStr).Error; err != nil {
		return err
	}

	// 更新缓存
	s.mutex.Lock()
	s.cache[key] = string(jsonValue) // 缓存中存储原始值
	s.mutex.Unlock()

	return nil
}

// Reload 强制重新加载配置
func (s *SysConfig) Reload() error {
	return s.loadConfig()
}

// NewSysConfig 创建系统配置实例
func NewSysConfig(db *gorm.DB) (*SysConfig, error) {
	s := &SysConfig{
		db:       db,
		cache:    make(map[string]string),
		cacheTTL: 5 * time.Minute, // 缓存5分钟
	}

	// 立即加载配置
	if err := s.loadConfig(); err != nil {
		return nil, fmt.Errorf("加载系统配置失败: %w", err)
	}

	return s, nil
}

// Init 初始化全局配置实例
func Init(db *gorm.DB) error {
	var err error
	once.Do(func() {
		var s *SysConfig
		s, err = NewSysConfig(db)
		if err != nil {
			return
		}
		instance = s
	})
	return err
}

// GetInstance 获取全局配置实例
func GetInstance() *SysConfig {
	if instance == nil {
		panic("系统配置未初始化")
	}
	return instance
}
