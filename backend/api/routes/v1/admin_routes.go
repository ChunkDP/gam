package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理员相关路由
func RegisterAdminRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	base := services.NewBaseCRUDService[models.Admin](db)
	cache := services.NewCacheBaseService(base, "admin")
	log := services.NewLogBaseService(cache, "admin")
	adminService := services.NewAdminService(db, log)

	h := handlers.NewAdminHandler(adminService)
	admins := r.Group("/admins")
	{
		admins.GET("", h.GetAdminList)
		admins.GET("/:id", h.GetAdmin)
		admins.POST("", h.CreateAdmin)
		admins.PUT("/:id", h.UpdateAdmin)
		admins.DELETE("/:id", h.DeleteAdmin)
		admins.PUT("/:id/status", h.UpdateAdminStatus)
		admins.PUT("/:id/password", h.UpdatePassword)
		admins.GET("/check-field", h.CheckAdminFieldUnique)
	}
}
