package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes 注册角色相关路由
func RegisterRoleRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	base := services.NewBaseCRUDService[models.Role](db)
	cache := services.NewCacheBaseService(base, "role")
	log := services.NewLogBaseService(cache, "role")
	roleService := services.NewRoleService(db, log)
	h := handlers.NewRoleHandler(roleService)
	// 角色管理
	roles := r.Group("/roles")
	{
		roles.GET("", h.GetRoleList)
		roles.GET("/:id", h.GetRole)
		roles.POST("", h.CreateRole)
		roles.PUT("/:id", h.UpdateRole)
		roles.DELETE("/:id", h.DeleteRole)
		roles.PUT("/:id/status", h.UpdateRoleStatus)
		roles.PUT("/:id/sort", h.UpdateRoleSort)
		roles.GET("/check-field", h.CheckRoleFieldUnique)

		// 角色权限管理
		roles.GET("/permissions/:roleId/menus", h.GetRoleMenus)
		roles.PUT("/permissions/:roleId/menus", h.UpdateRoleMenus)
	}
}
