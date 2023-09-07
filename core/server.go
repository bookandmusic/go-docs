package core

import (
	"strconv"

	"github.com/bookandmusic/docs/global"
)

func RunServer() {
	r := InitRouters()

	// 启动服务器
	port := global.GVA_CONFIG.Server.ServerPort
	global.GVA_LOG.Info("server start: ", port)
	r.Run(":" + strconv.Itoa(port))
}
