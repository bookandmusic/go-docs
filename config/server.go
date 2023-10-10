package config

type Server struct {
	ServerName       string `mapstructure:"server-name" json:"server-name" ini:"server-name"`
	Keyword          string `mapstructure:"keyword" json:"keyword" ini:"keyword"`
	Description      string `mapstructure:"description" json:"description" ini:"description"`
	ServerPort       int    `mapstructure:"server-port" json:"server-port" ini:"server-port"`
	DbType           string `mapstructure:"db-type" json:"db-type" ini:"db-type"`
	SecretKey        string `mapstructure:"secret-key" json:"secret-key" ini:"secret-key"`
	Debug            bool   `mapstructure:"debug" json:"debug" ini:"debug"`
	WorkingDirectory string `mapstructure:"work-dir" json:"work-dir" ini:"work-dir"`
	UploadPath       string `mapstructure:"upload" json:"upload" ini:"upload"`
}
