package config

type Log struct {
	Level   string `mapstructure:"level" json:"level" ini:"level"`
	LogFile string `mapstructure:"logfile" json:"logfile" ini:"logfile"`
}
