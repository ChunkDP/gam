package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterMenuRoutes 注册菜单相关路由
func RegisterMenuRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	menuService := services.NewMenuService(db)
	h := handlers.NewMenuHandler(menuService)

	menus := r.Group("/menus")
	{
		menus.GET("/tree", h.GetMenuTree)
		menus.GET("", h.GetMenuList)
		menus.POST("", h.CreateMenu)
		menus.PUT("/:id", h.UpdateMenu)
		menus.DELETE("/:id", h.DeleteMenu)
		menus.PUT("/:id/status", h.UpdateMenuStatus)
		menus.PUT("/:id/sort", h.UpdateMenuSort)
		menus.PUT("/:id/hidden", h.UpdateMenuHidden)
	}
}
