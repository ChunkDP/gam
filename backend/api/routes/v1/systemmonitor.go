package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// 注册系统监控路由
func RegisterSystemMonitorRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	systemMonitorService := services.NewSystemMonitorService(db)
	systemMonitorHandler := handlers.NewSystemMonitorHandler(systemMonitorService)

	systemMonitorGroup := r.Group("/system/monitor")
	{
		systemMonitorGroup.POST("/collect", systemMonitorHandler.CollectSystemInfo)
		systemMonitorGroup.GET("/list", systemMonitorHandler.GetSystemMonitors)
		systemMonitorGroup.GET("/latest", systemMonitorHandler.GetLatestSystemMonitor)
	}
}
