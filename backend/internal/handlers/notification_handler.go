package handlers

import (
	"log"
	"net/http"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/internal/services"
	"normaladmin/backend/pkg/utils/response"
	"normaladmin/backend/pkg/websocket"
	"strconv"
	"time"

	"normaladmin/backend/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notificationService services.NotificationService
	notificationHub     *websocket.NotificationHub
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	hub := websocket.NewNotificationHub()
	go hub.Run()

	return &NotificationHandler{
		notificationService: notificationService,
		notificationHub:     hub,
	}
}

// CreateNotificationType 创建通知类型
// @Summary 创建通知类型
// @Description 创建新的通知类型
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param notification_type body models.NotificationType true "通知类型信息"
// @Success 200 {object} response.Response
// @Router /api/notifications/types [post]
func (h *NotificationHandler) CreateNotificationType(c *gin.Context) {
	var notificationType models.NotificationType
	if err := c.ShouldBindJSON(&notificationType); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.notificationService.CreateNotificationType(&notificationType); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建通知类型失败: "+err.Error())
		return
	}

	response.Success(c, notificationType)
}

// UpdateNotificationType 更新通知类型
// @Summary 更新通知类型
// @Description 更新指定ID的通知类型
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知类型ID"
// @Param notification_type body models.NotificationType true "通知类型信息"
// @Success 200 {object} response.Response
// @Router /api/notifications/types/{id} [put]
func (h *NotificationHandler) UpdateNotificationType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知类型ID")
		return
	}

	var notificationType models.NotificationType
	if err := c.ShouldBindJSON(&notificationType); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.notificationService.UpdateNotificationType(uint(id), &notificationType); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新通知类型失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteNotificationType 删除通知类型
// @Summary 删除通知类型
// @Description 删除指定ID的通知类型
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知类型ID"
// @Success 200 {object} response.Response
// @Router /api/notifications/types/{id} [delete]
func (h *NotificationHandler) DeleteNotificationType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知类型ID")
		return
	}

	if err := h.notificationService.DeleteNotificationType(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除通知类型失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetNotificationType 获取通知类型
// @Summary 获取通知类型
// @Description 获取指定ID的通知类型
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知类型ID"
// @Success 200 {object} response.Response{data=models.NotificationType}
// @Router /api/notifications/types/{id} [get]
func (h *NotificationHandler) GetNotificationType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知类型ID")
		return
	}

	notificationType, err := h.notificationService.GetNotificationType(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取通知类型失败: "+err.Error())
		return
	}

	response.Success(c, notificationType)
}

// GetAllNotificationTypes 获取所有通知类型
// @Summary 获取所有通知类型
// @Description 获取系统中所有的通知类型
// @Tags 通知管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.NotificationType}
// @Router /api/notifications/types [get]
func (h *NotificationHandler) GetAllNotificationTypes(c *gin.Context) {
	types, err := h.notificationService.GetAllNotificationTypes()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取通知类型列表失败: "+err.Error())
		return
	}

	response.Success(c, types)
}

// CreateNotification 创建通知
// @Summary 创建通知
// @Description 创建新的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param notification body models.Notification true "通知信息"
// @Success 200 {object} response.Response
// @Router /api/notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	// 如果提供了过期时间字符串，转换为时间类型
	if notification.ExpirationTime != "" {
		// 解析时间字符串
		expirationTime, err := time.ParseInLocation("2006-01-02 15:04:05", notification.ExpirationTime, time.Local)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "无效的过期时间格式")
			return
		}

		// 验证过期时间是否在当前时间之后
		if expirationTime.Before(time.Now()) {
			response.Error(c, http.StatusBadRequest, "过期时间不能早于当前时间")
			return
		}
	}

	// 设置默认状态为草稿
	notification.Status = 0

	if err := h.notificationService.CreateNotification(&notification); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建通知失败: "+err.Error())
		return
	}

	response.Success(c, notification)
}

// UpdateNotification 更新通知
// @Summary 更新通知
// @Description 更新指定ID的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Param notification body models.Notification true "通知信息"
// @Success 200 {object} response.Response
// @Router /api/notifications/{id} [put]
func (h *NotificationHandler) UpdateNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.notificationService.UpdateNotification(uint(id), &notification); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新通知失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteNotification 删除通知
// @Summary 删除通知
// @Description 删除指定ID的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := h.notificationService.DeleteNotification(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除通知失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetNotification 获取通知详情
// @Summary 获取通知详情
// @Description 获取指定ID的通知详情
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response{data=models.Notification}
// @Router /api/notifications/{id} [get]
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	notification, err := h.notificationService.GetNotification(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取通知详情失败: "+err.Error())
		return
	}

	response.Success(c, notification)
}

// GetNotifications 获取通知列表
// @Summary 获取通知列表
// @Description 获取通知列表，支持分页和筛选
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param title query string false "通知标题"
// @Param type_id query int false "通知类型ID"
// @Param level query int false "重要程度"
// @Param status query int false "状态"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} response.Response{data=[]models.Notification,total=int64}
// @Router /api/notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	var query models.NotificationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的查询参数: "+err.Error())
		return
	}

	notifications, total, err := h.notificationService.GetNotifications(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取通知列表失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"notifications": notifications, "total": total})
}

// PublishNotification 发布通知
// @Summary 发布通知
// @Description 发布指定ID的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/notifications/{id}/publish [post]
func (h *NotificationHandler) PublishNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	// 获取当前用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}

	username := jwt.GetUsername(c)

	if err := h.notificationService.PublishNotification(uint(id), userID, username, userType); err != nil {
		response.Error(c, http.StatusInternalServerError, "发布通知失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// RecallNotification 撤回通知
// @Summary 撤回通知
// @Description 撤回指定ID的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/notifications/{id}/recall [post]
func (h *NotificationHandler) RecallNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := h.notificationService.RecallNotification(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "撤回通知失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUserNotifications 获取用户通知列表
// @Summary 获取用户通知列表
// @Description 获取当前用户的通知列表，支持分页和筛选
// @Tags 通知中心
// @Accept json
// @Produce json
// @Param is_read query bool false "是否已读"
// @Param type_id query int false "通知类型ID"
// @Param level query int false "重要程度"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} response.Response{data=[]models.NotificationReceiver,total=int64}
// @Router /api/user/notifications [get]
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	// 从JWT中获取用户信息

	// 获取当前用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}

	var query models.UserNotificationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的查询参数")
		return
	}

	notifications, total, err := h.notificationService.GetUserNotifications(userID, userType, &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取通知列表失败")
		return
	}

	response.Success(c, gin.H{
		"total": total,
		"items": notifications,
	})
}

// MarkNotificationAsRead 标记通知为已读
// @Summary 标记通知为已读
// @Description 标记指定ID的通知为已读
// @Tags 通知中心
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/user/notifications/{id}/read [post]
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	// 获取当前用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}

	if err := h.notificationService.MarkNotificationAsRead(uint(notificationID), userID, userType); err != nil {
		response.Error(c, http.StatusInternalServerError, "标记已读失败")
		return
	}

	response.Success(c, nil)
}

// MarkAllNotificationsAsRead 标记所有通知为已读
// @Summary 标记所有通知为已读
// @Description 标记当前用户的所有通知为已读
// @Tags 通知中心
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/user/notifications/read-all [post]
func (h *NotificationHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	// 获取当前用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}

	if err := h.notificationService.MarkAllNotificationsAsRead(userID, userType); err != nil {
		response.Error(c, http.StatusInternalServerError, "标记全部已读失败")
		return
	}

	response.Success(c, nil)
}

// DeleteUserNotification 删除用户通知
// @Summary 删除用户通知
// @Description 删除当前用户的指定通知
// @Tags 通知中心
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/user/notifications/{id} [delete]
func (h *NotificationHandler) DeleteUserNotification(c *gin.Context) {
	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	// 获取当前用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}

	if err := h.notificationService.DeleteUserNotification(uint(notificationID), userID, userType); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除通知失败")
		return
	}

	response.Success(c, nil)
}

// GetUnreadNotificationCount 获取未读通知数量
// @Summary 获取未读通知数量
// @Description 获取当前用户的未读通知数量
// @Tags 通知中心
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=int64}
// @Router /api/user/notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadNotificationCount(c *gin.Context) {
	// 获取当前用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}
	count, err := h.notificationService.GetUnreadNotificationCount(userID, userType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取未读通知数量失败")
		return
	}

	response.Success(c, gin.H{
		"count": count,
	})
}

// WebSocketHandler WebSocket连接处理
// @Summary WebSocket连接
// @Description 建立WebSocket连接以接收实时通知
// @Tags 通知中心
// @Accept json
// @Produce json
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/notifications [get]
func (h *NotificationHandler) WebSocketHandler(c *gin.Context) {
	// 获取用户信息
	userID := jwt.GetUserID(c)
	if userID == 0 {
		return // GetUserID 已经处理了错误响应
	}

	userType := jwt.GetUserType(c)
	if userType == "" {
		return // GetUserType 已经处理了错误响应
	}

	username := jwt.GetUsername(c)

	// 升级HTTP连接为WebSocket连接
	upgrader := gorillaws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		response.Error(c, http.StatusInternalServerError, "WebSocket连接建立失败")
		return
	}

	// 创建新的客户端连接
	client := websocket.NewClient(
		conn,
		userID,
		userType,
		username,
		h.notificationHub,
	)

	// 注册客户端
	h.notificationHub.Register <- client

	// 开始处理WebSocket消息
	go client.ReadPump()
	go client.WritePump()
}
