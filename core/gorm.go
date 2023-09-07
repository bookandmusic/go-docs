package core

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type DBBASE interface {
	GetLogMode() string
}

// Config gorm 自定义配置
func InitGormConfig(prefix string, singular bool) *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	_default := logger.New(log.New(global.GVA_LOG.Out, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
	})
	var logMode DBBASE
	switch global.GVA_CONFIG.Server.DbType {
	case "mysql":
		logMode = &global.GVA_CONFIG.MySQL
	case "sqlite":
		logMode = &global.GVA_CONFIG.SQLite
	default:
		logMode = &global.GVA_CONFIG.SQLite
	}

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}

func InitDatabase() *gorm.DB {
	switch global.GVA_CONFIG.Server.DbType {
	case "mysql":
		return GormMySQL()
	case "sqlite":
		return GormSqlite()
	default:
		return GormSqlite()
	}
}

func MigrateModels() {
	// 迁移所有模型
	db := global.GVA_DB
	err := db.AutoMigrate(
		models.User{},
		models.Category{},
		models.Tag{},
		models.Collection{},
		models.Article{},
		models.Setting{},
		models.Journal{},
	)
	if err != nil {
		panic(fmt.Errorf("register table failed: %v", err))
	}
	global.GVA_LOG.Info("register table success")
}
