package v1

import (
	"normaladmin/backend/database"
	"normaladmin/backend/internal/handlers"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/rabbitmq"
	"normaladmin/backend/pkg/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(r *gin.RouterGroup, rabbitmq *rabbitmq.RabbitMQ, notificationHub *websocket.NotificationHub) {

	db := database.GetDB()

	notificationService := services.NewNotificationService(db, rabbitmq, notificationHub)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// 通知类型管理
	notificationTypeGroup := r.Group("/notifications/types")
	{
		notificationTypeGroup.GET("", notificationHandler.GetAllNotificationTypes)
		notificationTypeGroup.GET("/:id", notificationHandler.GetNotificationType)
		notificationTypeGroup.POST("", notificationHandler.CreateNotificationType)
		notificationTypeGroup.PUT("/:id", notificationHandler.UpdateNotificationType)
		notificationTypeGroup.DELETE("/:id", notificationHandler.DeleteNotificationType)
	}

	// 通知管理
	notificationGroup := r.Group("/notifications")
	{
		notificationGroup.GET("", notificationHandler.GetNotifications)
		notificationGroup.GET("/:id", notificationHandler.GetNotification)
		notificationGroup.POST("", notificationHandler.CreateNotification)
		notificationGroup.PUT("/:id", notificationHandler.UpdateNotification)
		notificationGroup.DELETE("/:id", notificationHandler.DeleteNotification)
		notificationGroup.POST("/:id/publish", notificationHandler.PublishNotification)
		notificationGroup.POST("/:id/recall", notificationHandler.RecallNotification)
		notificationGroup.GET("/:id/stats", notificationHandler.GetNotificationStats)
	}

	// 用户通知
	userNotificationGroup := r.Group("/user/notifications")
	{
		userNotificationGroup.GET("", notificationHandler.GetUserNotifications)
		userNotificationGroup.POST("/:id/read", notificationHandler.MarkNotificationAsRead)
		userNotificationGroup.POST("/read-all", notificationHandler.MarkAllNotificationsAsRead)
		userNotificationGroup.DELETE("/:id", notificationHandler.DeleteUserNotification)
		userNotificationGroup.GET("/unread-count", notificationHandler.GetUnreadNotificationCount)
	}

}
