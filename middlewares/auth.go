package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求的路径
		path := c.FullPath()

		// 检查当前请求的路径是否以 "/admin/" 开头
		if strings.HasPrefix(path, "/admin/") && path != "/admin/login" {
			// 检查用户是否已登录
			session := sessions.Default(c)
			isLogged := session.Get("is_logged_in")

			if isLogged != true {
				c.Redirect(http.StatusFound, "/admin/login")
				c.Abort() // 终止请求处理
				return
			}
		}
		c.Next()
	}
}
