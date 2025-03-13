package response

import (
	"normaladmin/backend/pkg/logger"
	"runtime"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// response 内部使用的基础响应函数
func response(c *gin.Context, status int, data interface{}, errMsg string) {
	if errMsg != "" {
		_, file, line, ok := runtime.Caller(2) // 注意这里改为2，因为多了一层调用
		if ok {
			logger.Error(errMsg, logger.Field("file", file), logger.Field("line", line))
		} else {
			logger.Error(errMsg)
		}
		c.JSON(status, ResponseData{Error: errMsg})
	} else {
		c.JSON(status, ResponseData{Data: data})
	}
}

// Success 成功响应

func Success(c *gin.Context, data interface{}) {
	response(c, 200, data, "")
}

// Error 错误响应

func Error(c *gin.Context, status int, errMsg string) {
	response(c, status, nil, errMsg)
}
