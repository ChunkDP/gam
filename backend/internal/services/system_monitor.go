package services

import (
	"fmt"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/utils"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"gorm.io/gorm"
)

type SystemMonitorService interface {
	CollectSystemInfo() (*models.SystemMonitor, error)
	GetSystemMonitors(query *models.SystemMonitorQuery) ([]models.SystemMonitor, error)
	GetLatestSystemMonitor() (*models.SystemMonitor, error)
}

type systemMonitorService struct {
	db *gorm.DB
}

func NewSystemMonitorService(db *gorm.DB) SystemMonitorService {
	return &systemMonitorService{db: db}
}

// CollectSystemInfo 收集系统信息
func (s *systemMonitorService) CollectSystemInfo() (*models.SystemMonitor, error) {
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 获取内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 获取磁盘使用率
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// 获取网络IO
	netIOBefore, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	netIOAfter, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	var netIO string
	if len(netIOBefore) > 0 && len(netIOAfter) > 0 {
		bytesRecv := netIOAfter[0].BytesRecv - netIOBefore[0].BytesRecv
		bytesSent := netIOAfter[0].BytesSent - netIOBefore[0].BytesSent
		netIO = fmt.Sprintf("接收: %s/s, 发送: %s/s",
			utils.FormatBytes(float64(bytesRecv)),
			utils.FormatBytes(float64(bytesSent)))
	}

	// 获取进程数
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	// 获取负载均衡
	loadInfo, err := load.Avg()
	if err != nil && runtime.GOOS != "windows" {
		return nil, err
	}

	var loadAvg string
	if runtime.GOOS != "windows" {
		loadAvg = fmt.Sprintf("1分钟: %.2f, 5分钟: %.2f, 15分钟: %.2f",
			loadInfo.Load1, loadInfo.Load5, loadInfo.Load15)
	} else {
		loadAvg = "Windows系统不支持负载均衡指标"
	}

	// 创建系统监控记录
	monitor := &models.SystemMonitor{
		CPUUsage:     cpuPercent[0],
		MemoryUsage:  memInfo.UsedPercent,
		DiskUsage:    diskInfo.UsedPercent,
		NetworkIO:    netIO,
		ProcessCount: len(processes),
		LoadAverage:  loadAvg,
		CreatedAt:    time.Now(),
	}

	// 保存到数据库
	if err := s.db.Create(monitor).Error; err != nil {
		return nil, err
	}

	return monitor, nil
}

// GetSystemMonitors 获取系统监控列表
func (s *systemMonitorService) GetSystemMonitors(query *models.SystemMonitorQuery) ([]models.SystemMonitor, error) {
	var monitors []models.SystemMonitor
	db := s.db.Model(&models.SystemMonitor{})

	if query.StartTime != "" {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	limit := 100 // 默认限制
	if query.Limit > 0 {
		limit = query.Limit
	}

	err := db.Order("created_at DESC").Limit(limit).Find(&monitors).Error
	return monitors, err
}

// GetLatestSystemMonitor 获取最新的系统监控信息
func (s *systemMonitorService) GetLatestSystemMonitor() (*models.SystemMonitor, error) {
	var monitor models.SystemMonitor
	err := s.db.Order("created_at DESC").First(&monitor).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}
