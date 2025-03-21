package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterSystemRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	systemService := services.NewSystemService(db)
	h := handlers.NewSystemHandler(systemService)

	system := r.Group("/system")
	{
		// 日志管理
		system.GET("/logs", h.GetSystemLogs)
		system.DELETE("/logs", h.DeleteSystemLogs)

		// 系统监控
		system.GET("/monitor", h.GetSystemMonitor)
		system.GET("/monitor/history", h.GetMonitorHistory)

	}
}
