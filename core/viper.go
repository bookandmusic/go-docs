package core

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

func InitConfig(configFile string, serverPort int) *viper.Viper {
	v := viper.New()

	// 加载配置文件
	v.SetConfigFile(configFile)
	v.SetConfigType("ini")
	v.Set("mapstructure", "true")
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Fatal error config file: %s \n", err)
	}
	secretKey := utils.GenerateRandomKey(32)
	v.SetDefault("server.server-port", serverPort)
	v.SetDefault("server.db-type", "sqlite")
	v.SetDefault("server.secret-key", secretKey)
	v.SetDefault("server.debug", false)
	v.SetDefault("server.work-dir", "./")

	v.SetDefault("server.server-name", "GDocs")
	v.SetDefault("server.keyword", "layuimini,layui,gin,gdocs")
	v.SetDefault("server.description", "面向个人的文档和博客系统")

	v.SetDefault("log.level", "info")
	v.SetDefault("mysql.engine", "InnoDB")
	v.SetDefault("sqlite.log-mode", "info")
	v.SetDefault("mysql.log-mode", "info")

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			panic(fmt.Errorf("Failed to init config: %v", err))
		}
	})

	// 解析配置到 AppConfigInstance
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		panic(fmt.Errorf("Failed to init config: %v", err))
	}

	global.GVA_CONFIG.Server.UploadPath = filepath.Join(global.GVA_CONFIG.Server.WorkingDirectory, "uploads")
	global.GVA_CONFIG.SQLite.DBPath = filepath.Join(global.GVA_CONFIG.Server.WorkingDirectory, "database", "db.sqlite")
	global.GVA_CONFIG.Log.LogFile = filepath.Join(global.GVA_CONFIG.Server.WorkingDirectory, "logs/app.log")

	if err := utils.PathExists(global.GVA_CONFIG.Server.UploadPath); err != nil {
		panic(fmt.Errorf("Fatal error create upload path: %s \n", err))
	}
	return v
}
