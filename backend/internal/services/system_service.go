package services

import (
	"fmt"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/utils"
	"runtime"
	"time"

	"gorm.io/gorm"
)

type SystemService interface {
	// 日志管理
	CreateLog(log *models.SystemLog) error
	GetLogList(query map[string]interface{}, page, pageSize string) ([]models.SystemLog, int64, error)
	DeleteLogs(before time.Time) error

	// 系统监控
	CollectSystemInfo() (*models.SystemMonitor, error)
	GetMonitorData(duration string) ([]models.SystemMonitor, error)
}

type systemService struct {
	db *gorm.DB
}

func NewSystemService(db *gorm.DB) SystemService {
	return &systemService{db: db}
}

// 系统监控实现
func (s *systemService) CollectSystemInfo() (*models.SystemMonitor, error) {
	monitor := &models.SystemMonitor{
		CPUUsage:     utils.GetCPUUsage(),
		MemoryUsage:  utils.GetMemoryUsage(),
		DiskUsage:    utils.GetDiskUsage(),
		NetworkIO:    utils.GetNetworkIO(),
		ProcessCount: runtime.NumGoroutine(),
		LoadAverage:  utils.GetLoadAverage(),
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(monitor).Error; err != nil {
		return nil, err
	}

	return monitor, nil
}

func (s *systemService) DeleteLogs(before time.Time) error {
	return s.db.Where("created_at < ?", before).Delete(&models.SystemLog{}).Error
}

func (s *systemService) GetMonitorData(duration string) ([]models.SystemMonitor, error) {
	var monitors []models.SystemMonitor
	durationTime, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	startTime := time.Now().Add(-durationTime)
	err = s.db.Where("created_at >= ?", startTime).Find(&monitors).Error
	return monitors, err
}

func (s *systemService) CreateLog(log *models.SystemLog) error {
	// 设置默认值
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	// 验证必要字段
	if log.Module == "" {
		return fmt.Errorf("module is required")
	}
	if log.Action == "" {
		return fmt.Errorf("action is required")
	}
	if log.Method == "" {
		return fmt.Errorf("method is required")
	}
	if log.URL == "" {
		return fmt.Errorf("url is required")
	}
	if log.IP == "" {
		return fmt.Errorf("ip is required")
	}

	// 插入日志
	return s.db.Create(log).Error
}

func (s *systemService) GetLogList(query map[string]interface{}, page, pageSize string) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64
	db := s.db.Model(&models.SystemLog{})

	// 处理查询条件
	for field, value := range query {
		if strValue, ok := value.(string); ok && strValue != "" {
			db = db.Where(field+" = ?", strValue)
		} else if value != nil {
			db = db.Where(field+" = ?", value)
		}
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := db.Scopes(models.Paginate(page, pageSize)).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
