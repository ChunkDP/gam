package models

import (
	"time"
)

// SystemLog 系统日志
type SystemLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id"`                    // 操作用户ID
	Username  string    `json:"username"`                   // 操作用户名
	Module    string    `json:"module"`                     // 操作模块
	Action    string    `json:"action"`                     // 操作动作
	Method    string    `json:"method"`                     // 请求方法
	URL       string    `json:"url"`                        // 请求URL
	IP        string    `json:"ip"`                         // 请求IP
	UserAgent string    `json:"user_agent" gorm:"size:500"` // 用户代理
	Params    string    `json:"params" gorm:"type:text"`    // 请求参数
	Result    string    `json:"result" gorm:"type:text"`    // 操作结果
	Status    int       `json:"status"`                     // 状态码
	Duration  int64     `json:"duration"`                   // 执行时长(ms)
	CreatedAt time.Time `json:"created_at"`
}

// SystemMonitor 系统监控
type SystemMonitor struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	CPUUsage     float64   `json:"cpu_usage"`     // CPU使用率
	MemoryUsage  float64   `json:"memory_usage"`  // 内存使用率
	DiskUsage    float64   `json:"disk_usage"`    // 磁盘使用率
	NetworkIO    string    `json:"network_io"`    // 网络IO
	ProcessCount int       `json:"process_count"` // 进程数
	LoadAverage  string    `json:"load_average"`  // 负载均衡
	CreatedAt    time.Time `json:"created_at"`
}

// SystemMonitorQuery 系统监控查询参数
type SystemMonitorQuery struct {
	StartTime string `form:"start_time"` // 开始时间
	EndTime   string `form:"end_time"`   // 结束时间
	Limit     int    `form:"limit"`      // 限制条数
}
