package core

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

func InitConfig(configFile string) *viper.Viper {
	v := viper.New()

	// 加载配置文件
	v.SetConfigFile(configFile)
	v.SetConfigType("ini")

	v.Set("mapstructure", "true")
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	secretKey := utils.GenerateRandomKey(32)
	v.SetDefault("server-port", 8080)
	v.SetDefault("debug", false)
	v.SetDefault("work-dir", "./")
	v.SetDefault("server-name", "GDocs")
	v.SetDefault("secret-key", secretKey)
	v.SetDefault("keyword", "layuimini,layui,gin,gdocs")
	v.SetDefault("description", "面向个人的文档和博客系统")
	v.SetDefault("db-type", "sqlite")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.logfile", "./log/app.log")
	v.SetDefault("mysql.engine", "InnoDB")

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
	if global.GVA_CONFIG.SQLite.DBPath == "" {
		global.GVA_CONFIG.SQLite.DBPath = filepath.Join(global.GVA_CONFIG.Server.WorkingDirectory, "database", "db.sqlite")
	}

	if err := utils.PathExists(global.GVA_CONFIG.Server.UploadPath); err != nil {
		panic(fmt.Errorf("Fatal error create upload path: %s \n", err))
	}
	return v
}
