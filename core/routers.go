package core

import (
	"net/http"
	"time"

	_ "github.com/flosch/pongo2/v6"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	pongo2gin "gitlab.com/go-box/pongo2gin/v6"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/middlewares"
	routes "github.com/bookandmusic/docs/routers"
)

func InitRouters() *gin.Engine {
	if global.GVA_CONFIG.Server.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	// 使用 gin.Logger() 中间件记录请求信息到日志
	// 设置gin的日志输出为logrus
	r.Use(ginrus.Ginrus(global.GVA_LOG, time.RFC3339, true))

	// 配置Pongo2为模板引擎
	r.HTMLRender = pongo2gin.Default()

	r.NoRoute(middlewares.GlobalNotFoundMiddleware())

	// 指定静态文件目录
	r.Static("/static", "./static")
	// 设置上传的资源目录，将 "uploads" 目录映射为 URL "/uploads"
	r.StaticFS("/uploads", http.Dir(global.GVA_CONFIG.Server.UploadPath))

	store := cookie.NewStore([]byte(global.GVA_CONFIG.Server.SecretKey))
	r.Use(sessions.Sessions("docs", store))
	r.Use(middlewares.CSRFMiddleware())
	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.GlobalRecoveryMiddleware())
	// 使用 MinifyMiddleware 中间件
	r.Use(middlewares.MinifyMiddleware(global.GVA_MINIFY))

	routes.InitAdminRoutes(r)
	routes.InitBlogRoutes(r)

	return r
}
