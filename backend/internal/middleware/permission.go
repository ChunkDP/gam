package middleware

import (
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户信息
		// roleID, exists := c.Get("role_id")

		// if !exists {
		// 	response.Error(c, http.StatusUnauthorized, "Unauthorized")
		// 	c.Abort()
		// 	return
		// }

		// // 获取请求方法和路径
		// obj := c.Request.URL.Path
		// act := c.Request.Method
		// sub := fmt.Sprintf("%d", roleID)

		// // 检查权限
		// enforcer := auth.GetEnforcer()
		// ok, err := enforcer.Enforce(sub, obj, act)

		// fmt.Println("CasbinMiddleware", sub, obj, act)

		// if err != nil {
		// 	response.Error(c, http.StatusInternalServerError, "Permission check error")
		// 	c.Abort()
		// 	return
		// }

		// if !ok {
		// 	response.Error(c, http.StatusForbidden, "Permission denied")
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
