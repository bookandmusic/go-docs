package config

import (
	"fmt"
)

type SQLite struct {
	DBPath       string `mapstructure:"db-path" json:"db-path" ini:"db-path"`                      // 用于SQLite
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" ini:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" ini:"max-open-conns"` // 打开到数据库的最大连接数
	Prefix       string `mapstructure:"prefix" json:"prefix" ini:"prefix"`
	Singular     bool   `mapstructure:"singular" json:"singular" ini:"singular"` // 是否开启全局禁用复数，true表示开启
	LogMode      string `mapstructure:"log-mode" json:"log-mode" ini:"log-mode"` // 是否开启Gorm全局日志
}

func (s *SQLite) Dsn() string {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc", s.DBPath)
}

func (s *SQLite) GetLogMode() string {
	return s.LogMode
}
