package services

import (
	"encoding/json"
	"fmt"
	"log"
	"normaladmin/backend/pkg/rabbitmq"
	"normaladmin/backend/pkg/websocket"
)

// NotificationConsumer 通知消费者服务
type NotificationConsumer struct {
	rabbitmq        *rabbitmq.RabbitMQ
	notificationHub *websocket.NotificationHub
}

// NewNotificationConsumer 创建通知消费者服务
func NewNotificationConsumer(rabbitmq *rabbitmq.RabbitMQ, hub *websocket.NotificationHub) *NotificationConsumer {
	return &NotificationConsumer{
		rabbitmq:        rabbitmq,
		notificationHub: hub,
	}
}

// Start 启动消费者服务
func (nc *NotificationConsumer) Start() error {
	return nc.rabbitmq.ConsumeMessages("notifications", nc.handleNotification)
}

// handleNotification 处理通知消息
func (nc *NotificationConsumer) handleNotification(data []byte) error {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	fmt.Println("msg:", msg)
	// 根据消息类型处理
	switch msg["action"].(string) {
	case "new":
		// 通过WebSocket发送给在线用户
		if err := nc.notificationHub.SendToUser(
			uint(msg["user_id"].(float64)),
			msg["user_type"].(string),
			msg,
		); err != nil {
			log.Printf("发送WebSocket消息失败: %v", err)
		}

		// TODO: 对于离线用户，可以实现其他通知方式
		// 例如：发送邮件、推送等

	case "recall":
		notificationID := uint(msg["id"].(float64))

		// 1. 通过WebSocket通知在线用户
		recallMsg := map[string]interface{}{
			"type":    "notification",
			"action":  "recall",
			"id":      notificationID,
			"message": msg["message"],
		}

		// 广播撤回消息给所有相关用户
		if err := nc.notificationHub.Broadcast(recallMsg); err != nil {
			log.Printf("广播撤回消息失败: %v", err)
		}

		return nil
	default:
		return fmt.Errorf("未知的消息类型: %s", msg["action"])
	}

	return nil
}
