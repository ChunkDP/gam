package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// NotificationHub 管理WebSocket连接和消息广播
type NotificationHub struct {
	// 注册的客户端
	clients map[string]*Client
	// 注册请求
	Register chan *Client
	// 注销请求
	Unregister chan *Client
	// 广播消息
	broadcast chan []byte
	// 发送给特定用户的消息
	send chan *UserMessage
	// 互斥锁
	mutex sync.Mutex
}

// UserMessage 发送给特定用户的消息
type UserMessage struct {
	UserID   uint
	UserType string
	Data     []byte
}

// NewNotificationHub 创建一个新的通知中心
func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		send:       make(chan *UserMessage),
	}
}

// Run 启动通知中心
func (h *NotificationHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.clients[client.ID] = client
			h.mutex.Unlock()
			log.Printf("Client registered: %s (UserID: %d, UserType: %s)", client.ID, client.UserID, client.UserType)
		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mutex.Unlock()
			log.Printf("Client unregistered: %s", client.ID)
		case message := <-h.broadcast:
			h.mutex.Lock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
			h.mutex.Unlock()
		case message := <-h.send:
			h.mutex.Lock()
			for _, client := range h.clients {
				if client.UserID == message.UserID && client.UserType == message.UserType {
					select {
					case client.Send <- message.Data:
					default:
						close(client.Send)
						delete(h.clients, client.ID)
					}
				}
			}
			h.mutex.Unlock()
		}
	}
}

// SendToUser 发送消息给特定用户
func (h *NotificationHub) SendToUser(userID uint, userType string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	h.send <- &UserMessage{
		UserID:   userID,
		UserType: userType,
		Data:     data,
	}
	return nil
}

// Broadcast 广播消息给所有用户
func (h *NotificationHub) Broadcast(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	h.broadcast <- data
	return nil
}
