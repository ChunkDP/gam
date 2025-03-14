package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/rabbitmq"
	"normaladmin/backend/pkg/websocket"
	"time"

	"gorm.io/gorm"
)

// NotificationService 通知服务接口
type NotificationService interface {
	// 通知类型相关
	CreateNotificationType(notificationType *models.NotificationType) error
	UpdateNotificationType(id uint, notificationType *models.NotificationType) error
	DeleteNotificationType(id uint) error
	GetNotificationType(id uint) (*models.NotificationType, error)
	GetAllNotificationTypes() ([]models.NotificationType, error)

	// 通知相关
	CreateNotification(notification *models.Notification) error
	UpdateNotification(id uint, notification *models.Notification) error
	DeleteNotification(id uint) error
	GetNotification(id uint) (*models.Notification, error)
	GetNotifications(query *models.NotificationQuery) ([]models.Notification, int64, error)
	PublishNotification(id uint, currentUserID uint, currentUserName string, currentUserType string) error
	RecallNotification(id uint) error

	// 用户通知相关
	GetNotificationStats(notificationID uint) (map[string]int64, error)
	GetUserNotifications(userID uint, userType string, query *models.UserNotificationQuery) ([]models.NotificationReceiver, int64, error)
	MarkNotificationAsRead(notificationID uint, userID uint, userType string) error
	MarkAllNotificationsAsRead(userID uint, userType string) error
	DeleteUserNotification(notificationID uint, userID uint, userType string) error
	GetUnreadNotificationCount(userID uint, userType string) (int64, error)
	GetNotificationReceivers(notificationID uint, query *models.ReceiverQuery) ([]models.NotificationReceiver, int64, error)
}

type notificationService struct {
	db              *gorm.DB
	rabbitmq        *rabbitmq.RabbitMQ
	notificationHub *websocket.NotificationHub
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB, rabbitmq *rabbitmq.RabbitMQ, notificationHub *websocket.NotificationHub) NotificationService {
	return &notificationService{
		db:              db,
		rabbitmq:        rabbitmq,
		notificationHub: notificationHub,
	}
}

// CreateNotificationType 创建通知类型
func (s *notificationService) CreateNotificationType(notificationType *models.NotificationType) error {
	return s.db.Create(notificationType).Error
}

// UpdateNotificationType 更新通知类型
func (s *notificationService) UpdateNotificationType(id uint, notificationType *models.NotificationType) error {
	return s.db.Model(&models.NotificationType{}).Where("id = ?", id).Updates(notificationType).Error
}

// DeleteNotificationType 删除通知类型
func (s *notificationService) DeleteNotificationType(id uint) error {
	// 检查是否有通知使用此类型
	var count int64
	if err := s.db.Model(&models.Notification{}).Where("type_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该通知类型已被使用，无法删除")
	}
	return s.db.Delete(&models.NotificationType{}, id).Error
}

// GetNotificationType 获取通知类型
func (s *notificationService) GetNotificationType(id uint) (*models.NotificationType, error) {
	var notificationType models.NotificationType
	if err := s.db.First(&notificationType, id).Error; err != nil {
		return nil, err
	}
	return &notificationType, nil
}

// GetAllNotificationTypes 获取所有通知类型
func (s *notificationService) GetAllNotificationTypes() ([]models.NotificationType, error) {
	var types []models.NotificationType
	if err := s.db.Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// CreateNotification 创建通知
func (s *notificationService) CreateNotification(notification *models.Notification) error {
	// 如果提供了过期时间，转换为数据库格式
	if notification.ExpirationTime != "" {
		expirationTime, err := time.ParseInLocation("2006-01-02 15:04:05", notification.ExpirationTime, time.Local)
		if err != nil {
			return fmt.Errorf("无效的过期时间格式: %w", err)
		}
		// 直接设置解析后的时间到数据库字段
		notification.ExpirationTime = expirationTime.Format("2006-01-02 15:04:05")
	}

	return s.db.Create(notification).Error
}

// UpdateNotification 更新通知
func (s *notificationService) UpdateNotification(id uint, notification *models.Notification) error {
	// 检查通知状态，已发布的通知不能修改
	var existingNotification models.Notification
	if err := s.db.First(&existingNotification, id).Error; err != nil {
		return err
	}
	if existingNotification.Status == 1 {
		return errors.New("已发布的通知不能修改")
	}
	return s.db.Model(&models.Notification{}).Where("id = ?", id).Updates(notification).Error
}

// DeleteNotification 删除通知
func (s *notificationService) DeleteNotification(id uint) error {
	// 检查通知状态，已发布的通知不能删除
	var notification models.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return err
	}
	if notification.Status == 1 {
		return errors.New("已发布的通知不能删除")
	}
	return s.db.Delete(&models.Notification{}, id).Error
}

// GetNotification 获取通知详情
func (s *notificationService) GetNotification(id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := s.db.Preload("Type").First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetNotifications 获取通知列表
func (s *notificationService) GetNotifications(query *models.NotificationQuery) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	db := s.db.Model(&models.Notification{}).Preload("Type")

	// 应用查询条件
	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if query.TypeID != 0 {
		db = db.Where("type_id = ?", query.TypeID)
	}
	if query.Level != 0 {
		db = db.Where("level = ?", query.Level)
	}
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.StartTime != "" {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// PublishNotification 发布通知
func (s *notificationService) PublishNotification(id uint, currentUserID uint, currentUserName string, currentUserType string) error {
	// 获取通知详情
	var notification models.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return fmt.Errorf("获取通知失败: %w", err)
	}
	fmt.Println("notification:", notification)

	// 检查通知状态
	if notification.Status != 0 { // 0表示草稿状态
		return errors.New("只能发布草稿状态的通知")
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新通知状态
	now := time.Now()
	notification.Status = 1 // 1表示已发布
	notification.PublishTime = &now
	notification.SenderID = currentUserID
	notification.SenderName = currentUserName

	if err := tx.Save(&notification).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新通知状态失败: %w", err)
	}

	// 确定接收者并创建通知接收记录
	var userInfos []UserInfo
	switch notification.ReceiverType {
	case "all":
		// 获取所有会员
		if err := s.db.Model(&models.Member{}).
			Select("id as user_id, 'member' as user_type, username").
			Scan(&userInfos).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("获取会员列表失败: %w", err)
		}
		// 获取所有管理员
		var adminInfos []UserInfo
		if err := s.db.Model(&models.Admin{}).
			Select("id as user_id, 'admin' as user_type, username").
			Scan(&adminInfos).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("获取管理员列表失败: %w", err)
		}
		userInfos = append(userInfos, adminInfos...)

	case "members":
		// 只发送给会员
		if err := s.db.Model(&models.Member{}).
			Select("id as user_id, 'member' as user_type, username").
			Scan(&userInfos).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("获取会员列表失败: %w", err)
		}

	case "admins":
		// 只发送给管理员
		if err := s.db.Model(&models.Admin{}).
			Select("id as user_id, 'admin' as user_type, username").
			Scan(&userInfos).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("获取管理员列表失败: %w", err)
		}

	default:
		tx.Rollback()
		return errors.New("无效的接收者类型")
	}

	// 批量创建通知接收记录
	var receivers []models.NotificationReceiver
	for _, user := range userInfos {
		receivers = append(receivers, models.NotificationReceiver{
			NotificationID: notification.ID,
			UserID:         user.UserID,
			UserType:       user.UserType,
			UserName:       user.Username,
			IsRead:         false,
		})
	}

	if len(receivers) > 0 {
		if err := tx.CreateInBatches(receivers, 100).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建通知接收记录失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 发送WebSocket消息
	notificationMsg := map[string]interface{}{
		"type":       "notification",
		"action":     "new",
		"user_id":    currentUserID,
		"user_type":  currentUserType,
		"id":         notification.ID,
		"title":      notification.Title,
		"content":    notification.Content,
		"level":      notification.Level,
		"createTime": notification.CreatedAt,
	}
	// 将消息发送到RabbitMQ
	msgBytes, err := json.Marshal(notificationMsg)
	if err != nil {
		log.Printf("序列化通知消息失败: %v", err)
		return err
	}

	// 发送到RabbitMQ的通知队列
	if err := s.rabbitmq.PublishMessage("notifications", msgBytes); err != nil {
		log.Printf("发送RabbitMQ消息失败: %v", err)
		return err
	}

	return nil
}

// UserInfo 用户信息结构体
type UserInfo struct {
	UserID   uint   `gorm:"column:user_id"`
	UserType string `gorm:"column:user_type"`
	Username string `gorm:"column:username"`
}

// RecallNotification 撤回通知
func (s *notificationService) RecallNotification(id uint) error {
	// 获取通知详情
	var notification models.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return err
	}

	// 检查通知状态
	if notification.Status != 1 {
		return errors.New("只有已发布的通知可以撤回")
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新通知状态
	if err := tx.Model(&notification).Update("status", 2).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 更新所有接收记录为已撤回
	now := time.Now()
	// 更新所有接收记录的撤回状态
	if err := tx.Model(&models.NotificationReceiver{}).
		Where("notification_id = ?", id).
		Updates(map[string]interface{}{
			"is_recalled": true,
			"recall_time": now,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新通知接收记录失败: %w", err)
	}
	// 发送撤回消息到RabbitMQ
	recallMsg := map[string]interface{}{
		"type":    "notification",
		"action":  "recall",
		"id":      notification.ID,
		"message": "通知已被撤回",
	}
	msgBytes, err := json.Marshal(recallMsg)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.rabbitmq.PublishMessage("notification_recall", msgBytes); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetUserNotifications 获取用户通知列表
func (s *notificationService) GetUserNotifications(userID uint, userType string, query *models.UserNotificationQuery) ([]models.NotificationReceiver, int64, error) {
	var receivers []models.NotificationReceiver
	var total int64

	db := s.db.Model(&models.NotificationReceiver{}).
		Where("user_id = ? AND user_type = ? AND is_deleted = ?", userID, userType, false).
		Preload("Notification.Type")
	if query.ShowRecalled != nil {
		db = db.Where("is_recalled = ?", *query.ShowRecalled)
	}
	// 应用查询条件
	if query.IsRead != nil {
		db = db.Where("is_read = ?", *query.IsRead)
	}
	if query.TypeID != 0 {
		db = db.Joins("JOIN notifications ON notification_receivers.notification_id = notifications.id").
			Where("notifications.type_id = ?", query.TypeID)
	}
	if query.Level != 0 {
		db = db.Joins("JOIN notifications ON notification_receivers.notification_id = notifications.id").
			Where("notifications.level = ?", query.Level)
	}
	if query.StartTime != "" {
		db = db.Where("notification_receivers.created_at >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("notification_receivers.created_at <= ?", query.EndTime)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := db.Debug().Order("notification_receivers.created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&receivers).Error; err != nil {
		return nil, 0, err
	}

	return receivers, total, nil
}

// MarkNotificationAsRead 标记通知为已读
func (s *notificationService) MarkNotificationAsRead(notificationID uint, userID uint, userType string) error {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// 检查是否已读
	var receiver models.NotificationReceiver
	err := tx.Where("notification_id = ? AND user_id = ? AND user_type = ?",
		notificationID, userID, userType).First(&receiver).Error
	if err != nil {
		tx.Rollback()
		return errors.New("通知不存在或无权限访问")
	}

	// 如果未读，则更新已读状态并增加已读计数
	if !receiver.IsRead {
		// 更新接收记录为已读
		if err := tx.Model(&receiver).Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": now,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 增加通知的已读计数
		if err := tx.Model(&models.Notification{}).
			Where("id = ?", notificationID).
			UpdateColumn("read_count", gorm.Expr("read_count + ?", 1)).
			Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func (s *notificationService) MarkAllNotificationsAsRead(userID uint, userType string) error {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// 获取所有未读的通知接收记录
	var unreadReceivers []models.NotificationReceiver
	if err := tx.Where("user_id = ? AND user_type = ? AND is_read = ?",
		userID, userType, false).Find(&unreadReceivers).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 按通知ID分组，统计每个通知新增的已读数量
	notificationCounts := make(map[uint]int64)
	for _, receiver := range unreadReceivers {
		notificationCounts[receiver.NotificationID]++
	}

	// 更新接收记录为已读
	if err := tx.Model(&models.NotificationReceiver{}).
		Where("user_id = ? AND user_type = ? AND is_read = ?", userID, userType, false).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": now,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 批量更新通知的已读计数
	for notificationID, count := range notificationCounts {
		if err := tx.Model(&models.Notification{}).
			Where("id = ?", notificationID).
			UpdateColumn("read_count", gorm.Expr("read_count + ?", count)).
			Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteUserNotification 删除用户通知（软删除）
func (s *notificationService) DeleteUserNotification(notificationID uint, userID uint, userType string) error {
	result := s.db.Model(&models.NotificationReceiver{}).
		Where("notification_id = ? AND user_id = ? AND user_type = ?", notificationID, userID, userType).
		Update("is_deleted", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("通知不存在或无权限访问")
	}

	return nil
}

// GetUnreadNotificationCount 获取未读通知数量
func (s *notificationService) GetUnreadNotificationCount(userID uint, userType string) (int64, error) {
	var count int64
	err := s.db.Model(&models.NotificationReceiver{}).
		Where("user_id = ? AND user_type = ? AND is_read = ? AND is_deleted = ?", userID, userType, false, false).
		Count(&count).Error
	return count, err
}

// GetNotificationReceivers 获取通知接收者列表
func (s *notificationService) GetNotificationReceivers(notificationID uint, query *models.ReceiverQuery) ([]models.NotificationReceiver, int64, error) {
	var receivers []models.NotificationReceiver
	var total int64

	db := s.db.Model(&models.NotificationReceiver{}).
		Where("notification_id = ?", notificationID)

	// 应用查询条件
	if query.UserType != "" {
		db = db.Where("user_type = ?", query.UserType)
	}
	if query.IsRead != nil {
		db = db.Where("is_read = ?", *query.IsRead)
	}
	if query.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+query.UserName+"%")
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&receivers).Error; err != nil {
		return nil, 0, err
	}

	return receivers, total, nil
}

// GetNotificationStats 获取通知统计信息
func (s *notificationService) GetNotificationStats(notificationID uint) (map[string]int64, error) {
	var stats struct {
		TotalReceivers int64
		ReadCount      int64
		RecalledRead   int64 // 撤回前已读数量
	}

	// 获取总接收者数量
	if err := s.db.Model(&models.NotificationReceiver{}).
		Where("notification_id = ?", notificationID).
		Count(&stats.TotalReceivers).Error; err != nil {
		return nil, err
	}

	// 获取已读数量
	if err := s.db.Model(&models.NotificationReceiver{}).
		Where("notification_id = ? AND is_read = ?", notificationID, true).
		Count(&stats.ReadCount).Error; err != nil {
		return nil, err
	}

	// 获取撤回前已读数量
	if err := s.db.Model(&models.NotificationReceiver{}).
		Where("notification_id = ? AND is_read = ? AND read_time < recall_time",
			notificationID, true).
		Count(&stats.RecalledRead).Error; err != nil {
		return nil, err
	}

	return map[string]int64{
		"total_receivers": stats.TotalReceivers,
		"read_count":      stats.ReadCount,
		"recalled_read":   stats.RecalledRead,
	}, nil
}
