package config

type AppConfig struct {
	Server Server `mapstructure:"server" json:"server" ini:"server"`
	MySQL  MySQL  `mapstructure:"mysql" json:"mysql" ini:"mysql"`
	SQLite SQLite `mapstructure:"sqlite" json:"sqlite" ini:"sqlite"`
	Log    Log    `mapstructure:"log" json:"log" ini:"log"`
}
