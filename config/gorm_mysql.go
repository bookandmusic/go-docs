package config

import (
	"fmt"
)

type MySQL struct {
	Host         string `mapstructure:"host" json:"host" ini:"host"`
	Port         int    `mapstructure:"port" json:"port" ini:"port"`
	DbName       string `mapstructure:"db-name" json:"db-name" ini:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" ini:"username"`                   // 数据库密码
	Password     string `mapstructure:"password" json:"password" ini:"password"`                   // 数据库密码
	Engine       string `mapstructure:"engine" json:"engine" ini:"engine" default:"InnoDB"`        // 数据库引擎，默认InnoDB
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" ini:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" ini:"max-open-conns"` // 打开到数据库的最大连接数
	Prefix       string `mapstructure:"prefix" json:"prefix" ini:"prefix"`
	Singular     bool   `mapstructure:"singular" json:"singular" ini:"singular"` // 是否开启全局禁用复数，true表示开启
	LogMode      string `mapstructure:"log-mode" json:"log-mode" ini:"log-mode"` // 是否开启Gorm全局日志
}

func (m *MySQL) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.DbName)
}

func (m *MySQL) GetLogMode() string {
	return m.LogMode
}
