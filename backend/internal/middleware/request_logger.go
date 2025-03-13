package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResponseBodyWriter 自定义响应写入器，用于捕获响应内容
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法，同时写入原始响应和缓冲区
func (r ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// RequestLoggerMiddleware 请求日志中间件
func RequestLoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取请求信息
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 获取请求参数
		var params string
		if c.Request.Method == "GET" {
			// 对于 GET 请求，记录查询参数
			params = c.Request.URL.RawQuery
		} else {
			// 对于 POST/PUT 等请求，记录请求体
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				// 重新设置请求体，因为读取后 body 会被消费
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// 如果是 JSON 格式，美化输出
				if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
					var prettyJSON bytes.Buffer
					if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
						params = prettyJSON.String()
					} else {
						params = string(bodyBytes)
					}
				} else {
					params = string(bodyBytes)
				}
			}
		}

		// 获取用户信息
		var userID uint
		if userIDValue, exists := c.Get("user_id"); exists {
			if id, ok := userIDValue.(uint); ok {
				userID = id
			} else if idFloat, ok := userIDValue.(float64); ok {
				userID = uint(idFloat)
			} else if idInt, ok := userIDValue.(int); ok {
				userID = uint(idInt)
			}
		}

		// 如果需要用户名，可以根据 userID 查询数据库

		if userID > 0 {
			// 异步查询用户名，避免阻塞请求
			go func(uid uint) {
				var admin models.Admin
				if err := db.Select("username").First(&admin, uid).Error; err == nil {
					// 更新日志记录中的用户名
					db.Model(&models.SystemLog{}).Where("user_id = ? AND username = ''", uid).
						Updates(map[string]interface{}{"username": admin.Username})
				}
			}(userID)
		}
		// 包装响应写入器以捕获响应内容
		responseWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		duration := endTime.Sub(startTime).Milliseconds()

		// 获取响应状态和内容
		status := c.Writer.Status()
		responseBody := responseWriter.body.String()

		// 截断过长的响应内容
		if len(responseBody) > 1000 {
			responseBody = responseBody[:1000] + "... [截断]"
		}

		// 确定模块和操作
		module := getModuleFromPath(path)
		action := getActionFromMethod(method, path)

		// 创建日志记录
		log := models.SystemLog{
			UserID:    userID,
			Module:    module,
			Action:    action,
			Method:    method,
			URL:       path,
			IP:        ip,
			UserAgent: userAgent,
			Params:    params,
			Result:    responseBody,
			Status:    status,
			Duration:  duration,
			CreatedAt: endTime,
		}

		// 异步保存日志到数据库
		go func(log models.SystemLog) {
			// 如果有用户ID，尝试查询用户名
			if log.UserID > 0 {
				var admin models.Admin
				if err := db.Select("username").First(&admin, log.UserID).Error; err == nil {
					log.Username = admin.Username
				}
			}

			if err := db.Create(&log).Error; err != nil {
				// 如果数据库记录失败，则使用 zap 记录错误
				logger.Error("保存请求日志失败",
					logger.Field("error", err),
					logger.Field("log", fmt.Sprintf("%+v", log)),
				)
			}
		}(log)

		// 同时使用 zap 记录关键信息
		logLevel := "INFO"
		if status >= 400 {
			logLevel = "ERROR"
		}

		logMsg := fmt.Sprintf("[%s] %s %s - %d (%dms)", logLevel, method, path, status, duration)
		if status >= 400 {
			logger.Error(logMsg,
				logger.Field("ip", ip),
				logger.Field("user_id", userID),
				logger.Field("params", params),
				logger.Field("error", responseBody),
			)
		} else {
			logger.Info(logMsg,
				logger.Field("ip", ip),
				logger.Field("user_id", userID),
			)
		}
	}
}

// getModuleFromPath 从路径获取模块名
func getModuleFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return parts[2] // 例如 /api/admin/users -> admin
	}
	return "system"
}

// getActionFromMethod 从请求方法和路径获取操作类型
func getActionFromMethod(method, path string) string {
	switch method {
	case "GET":
		if strings.Contains(path, "/export") {
			return "导出"
		}
		return "查询"
	case "POST":
		return "创建"
	case "PUT":
		return "更新"
	case "DELETE":
		return "删除"
	default:
		return method
	}
}
