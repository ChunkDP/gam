package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// 配置相关路由
func RegisterConfigRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	configService := services.NewConfigService(db)
	h := handlers.NewConfigHandler(configService)
	config := r.Group("/configs")
	{
		config.GET("/groups", h.GetConfigGroups)
		config.GET("/items/:groupId", h.GetConfigItems)
		config.PUT("/value", h.UpdateConfigValue)
		config.PUT("/batch", h.BatchUpdateConfigs)
	}
}
