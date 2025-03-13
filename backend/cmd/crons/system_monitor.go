package crons

import (
	"log"
	"normaladmin/backend/internal/services"
	"time"

	"github.com/robfig/cron/v3"
)

// SetupSystemMonitorCron 设置系统监控定时任务
func SetupSystemMonitorCron(systemMonitorService services.SystemMonitorService) *cron.Cron {
	c := cron.New(cron.WithSeconds())

	// 每5分钟收集一次系统信息
	_, err := c.AddFunc("0 */5 * * * *", func() {
		_, err := systemMonitorService.CollectSystemInfo()
		if err != nil {
			log.Printf("收集系统信息失败: %v", err)
		} else {
			log.Printf("系统信息收集成功: %s", time.Now().Format("2006-01-02 15:04:05"))
		}
	})

	if err != nil {
		log.Fatalf("添加系统监控定时任务失败: %v", err)
	}

	c.Start()
	return c
}
