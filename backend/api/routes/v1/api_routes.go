package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册API管理相关路由
func RegisterAPIRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	base := services.NewBaseCRUDService[models.API](db)
	apiService := services.NewAPIService(db, base)
	h := handlers.NewAPIHandler(apiService)

	apis := r.Group("/apis")
	{
		apis.GET("", h.GetAPIList)
		apis.GET("/:id", h.GetAPI)
		apis.POST("", h.CreateAPI)
		apis.PUT("/:id", h.UpdateAPI)
		apis.DELETE("/:id", h.DeleteAPI)
		apis.POST("/:id/test", h.TestAPI)
	}
}
