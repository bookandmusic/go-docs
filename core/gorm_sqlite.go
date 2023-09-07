package core

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

func GormSqlite() *gorm.DB {
	s := global.GVA_CONFIG.SQLite

	err := utils.PathExists(s.DBPath)
	if err != nil {
		panic(fmt.Errorf("create sqlite file error: %v", err))
	}

	if db, err := gorm.Open(sqlite.Open(s.Dsn()), InitGormConfig(s.Prefix, s.Singular)); err != nil {
		panic(fmt.Errorf("open sqlite error: %v", err))
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		return db
	}
}
