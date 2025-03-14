package handlers

import (
	"net/http"
	"normaladmin/backend/internal/middleware"
	"normaladmin/backend/pkg/utils/response"
	"normaladmin/backend/pkg/websocket"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
)

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	notificationHub *websocket.NotificationHub
	secretKey       string
}

// NewWebSocketHandler 创建WebSocket处理器
func NewWebSocketHandler(notificationHub *websocket.NotificationHub, secretKey string) *WebSocketHandler {
	return &WebSocketHandler{
		notificationHub: notificationHub,
		secretKey:       secretKey,
	}
}

// handleWS 处理WebSocket连接的通用逻辑
func (h *WebSocketHandler) WebSocketHandler(c *gin.Context) {
	urltoken := c.Query("token")
	if urltoken == "" {
		response.Error(c, http.StatusUnauthorized, "未提供认证token")
		return
	}

	// 验证token并获取用户信息
	token, err := middleware.ParseToken(urltoken, h.secretKey)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "无效的token")
		return
	}

	upgrader := gorillaws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 可以根据需要设置更严格的来源检查
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	if claims, ok := token.Claims.(*middleware.Claims); ok && token.Valid {
		client := websocket.NewClient(
			conn,
			claims.UserID,
			claims.UserType,
			claims.Username,
			h.notificationHub,
		)

		h.notificationHub.Register <- client

		go client.ReadPump()
		go client.WritePump()
	} else {
		response.Error(c, http.StatusUnauthorized, "Invalid token claims")
		return
	}

}
