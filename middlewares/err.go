package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
)

func GlobalNotFoundMiddleware() gin.HandlerFunc {
	site_info := common.GenerateSiteInfo()
	return func(c *gin.Context) {
		// 渲染 404 模板
		err_msg := fmt.Sprintf("无法定位到该路径: %s", c.Request.URL.Path)
		c.HTML(http.StatusNotFound, "public/404.html", pongo2.Context{
			"err_msg":   err_msg, // 获取当前路由地址
			"site_info": site_info,
		})
	}
}

func GlobalRecoveryMiddleware() gin.HandlerFunc {
	site_info := common.GenerateSiteInfo()
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 在这里可以记录错误日志或其他处理
				path := c.Request.URL.Path
				global.GVA_LOG.Error(fmt.Sprintf("URL: %s, Server Error: %v", path, r))
				if strings.HasPrefix(path, common.APIPrefix) {
					c.JSON(http.StatusInternalServerError, common.ServerError)
					return
				}
				err_msg := fmt.Sprintf("无法定位到该路径: %s", c.Request.URL.Path)
				c.HTML(http.StatusNotFound, "public/404.html", pongo2.Context{
					"err_msg":   err_msg, // 获取当前路由地址
					"site_info": site_info,
				})
			}
		}()

		c.Next()
	}
}
