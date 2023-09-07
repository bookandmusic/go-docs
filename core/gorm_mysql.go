package core

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
)

func GormMySQL() *gorm.DB {
	m := global.GVA_CONFIG.MySQL
	if m.DbName == "" {
		panic(fmt.Errorf("mysql db_name not exists"))
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), InitGormConfig(m.Prefix, m.Singular)); err != nil {
		panic(fmt.Errorf("open mysql error: %v", err))
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
