package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"normaladmin/backend/pkg/sysconfig"
	"normaladmin/backend/pkg/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadAuth 上传认证中间件
func UploadAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否需要登录
		if !isPublicPath(c.Request.URL.Path) {
			if c.IsAborted() {
				return
			}
		}
		c.Next()
	}
}

// AntiLeech 防盗链中间件
func AntiLeech() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("AntiLeech")
		// 获取 Referer
		referer := c.Request.Header.Get("Referer")
		fmt.Println("referer:", referer)
		if referer == "" {
			response.Error(c, http.StatusForbidden, "禁止访问")
			c.Abort()
			return
		}

		// 解析 Referer
		u, err := url.Parse(referer)
		if err != nil {
			response.Error(c, http.StatusForbidden, "非法请求")
			c.Abort()
			return
		}

		// 获取允许的域名列表
		allowedDomains := strings.Split(sysconfig.Get("upload_allowed_domains", ""), ",")
		fmt.Println("allowedDomains:", allowedDomains)
		isAllowed := false
		for _, domain := range allowedDomains {

			fmt.Println("domain:", domain, "u.Host:", u.Host)
			if strings.TrimSpace(domain) == u.Host {
				fmt.Println("isAllowed:", isAllowed)
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.Error(c, http.StatusForbidden, "非法来源")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 检查是否是公开路径
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/uploads/public/", // 公开文件目录
	}
	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
