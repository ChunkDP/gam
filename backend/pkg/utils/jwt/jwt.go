package jwt

import (
	"normaladmin/backend/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

const (
	// JWT claims 中的 key 常量
	UserIDKey   = "user_id"
	UserTypeKey = "user_type"
	UsernameKey = "username"
)

// GetUserID 从 gin.Context 获取用户ID
func GetUserID(c *gin.Context) uint {
	value, exists := c.Get(UserIDKey)
	if !exists {

		response.Error(c, 401, "未获取到用户ID")
		return 0
	}

	if id, ok := value.(uint); ok {

		return id
	}
	return 0
}

// GetUserType 从 gin.Context 获取用户类型
func GetUserType(c *gin.Context) string {
	value, exists := c.Get(UserTypeKey)
	if !exists {
		response.Error(c, 401, "未获取到用户类型")
		return ""
	}
	if userType, ok := value.(string); ok {
		return userType
	}
	return ""
}

// GetUsername 从 gin.Context 获取用户名
func GetUsername(c *gin.Context) string {
	value, exists := c.Get(UsernameKey)
	if !exists {
		return "anonymous"
	}
	if username, ok := value.(string); ok {
		return username
	}
	return "anonymous"
}
