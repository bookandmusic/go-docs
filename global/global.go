package global

import (
	bleve "github.com/blevesearch/bleve/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	minify "github.com/tdewolff/minify/v2"
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/config"
)

var (
	GVA_DB          *gorm.DB
	GVA_CONFIG      config.AppConfig
	GVA_VP          *viper.Viper
	GVA_LOG         *logrus.Logger
	GVA_MINIFY      *minify.M
	GVA_BLEVE_INDEX bleve.Index
)
