package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/middleware"
	"normaladmin/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(r *gin.RouterGroup) {
	db := database.GetDB()
	uploadService := services.NewUploadService(db)
	h := handlers.NewUploadHandler(uploadService)

	upload := r.Group("/upload")
	upload.Use(middleware.UploadAuth())

	{
		upload.POST("", h.UploadFile)            // 通用文件上传
		upload.POST("/image", h.UploadImage)     // 图片上传
		upload.GET("/config", h.GetUploadConfig) // 获取上传配置
		upload.DELETE("", h.DeleteFile)          // 删除文件
		upload.POST("/batch", h.BatchUpload)
	}

}
