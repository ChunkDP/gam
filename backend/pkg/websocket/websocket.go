package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second
	// 读取超时时间
	pongWait = 60 * time.Second
	// Ping周期，必须小于pongWait
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小
	maxMessageSize = 512
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket 封装websocket连接
type WebSocket struct {
	conn *gorillaws.Conn
}

// Client WebSocket客户端连接
type Client struct {
	ID       string           // 客户端唯一标识
	Conn     *gorillaws.Conn  // WebSocket连接
	UserID   uint             // 用户ID
	UserType string           // 用户类型
	Username string           // 用户名
	Hub      *NotificationHub // 通知中心
	Send     chan []byte      // 发送消息的通道
}

// NewWebSocket 创建新的WebSocket连接
func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocket{conn: conn}, nil
}

// ReadMessage 读取消息
func (ws *WebSocket) ReadMessage() ([]byte, error) {
	_, message, err := ws.conn.ReadMessage()
	return message, err
}

// WriteMessage 写入消息
func (ws *WebSocket) WriteMessage(message []byte) error {
	return ws.conn.WriteMessage(gorillaws.TextMessage, message)
}

// Close 关闭连接
func (ws *WebSocket) Close() error {
	return ws.conn.Close()
}

// ReadPump 持续读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if gorillaws.IsUnexpectedCloseError(err, gorillaws.CloseGoingAway, gorillaws.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// WritePump 持续向客户端写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道已关闭
				c.Conn.WriteMessage(gorillaws.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(gorillaws.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorillaws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// NewClient 创建新的客户端连接
func NewClient(conn *gorillaws.Conn, userID uint, userType string, username string, hub *NotificationHub) *Client {
	return &Client{
		ID:       fmt.Sprintf("%s_%d", userType, userID),
		Conn:     conn,
		UserID:   userID,
		UserType: userType,
		Username: username,
		Hub:      hub,
		Send:     make(chan []byte, 256), // 缓冲区大小设为256
	}
}
