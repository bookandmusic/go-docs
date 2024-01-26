package middlewares

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求的路径
		path := c.Request.URL.Path

		// 检查当前请求的路径是否以 "/admin/" 开头
		if strings.HasPrefix(path, common.APIPrefix) && !strings.HasSuffix(path, common.APILoginUrl) && c.Request.Method != "OPTIONS" {
			// 检查用户是否已登录
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.JSON(http.StatusOK, common.TokenInvalidError)
				c.Abort() // 终止请求处理
				return
			}
			claims := &common.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(global.GVA_CONFIG.Server.SecretKey), nil
			})
			if err != nil {
				c.JSON(http.StatusOK, common.TokenInvalidError)
				c.Abort() // 终止请求处理
				return
			}

			if !token.Valid {
				c.JSON(http.StatusOK, common.TokenInvalidError)
				c.Abort() // 终止请求处理
				return
			}

			// 将Claims存入请求上下文
			c.Set("claims", claims)
			c.Next()

		}
		c.Next()
	}
}
