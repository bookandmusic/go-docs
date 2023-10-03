package main

import (
	"flag"
	"os"

	cli "github.com/urfave/cli/v2"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/bookandmusic/docs/commands"
	"github.com/bookandmusic/docs/core"
	"github.com/bookandmusic/docs/global"
)

var (
	configFile string

	username string
	password string
)

func init() {
	flag.StringVar(&configFile, "c", "./config/config.ini", "path to config file")
	flag.Parse()
}

func main() {
	// 初始化配置
	global.GVA_VP = core.InitConfig(configFile)
	global.GVA_LOG = core.InitLog()
	// 初始化数据库连接
	global.GVA_DB = core.InitDatabase()
	global.GVA_MINIFY = core.InitMinify()

	if fileLogger, ok := global.GVA_LOG.Out.(*lumberjack.Logger); ok {
		// 在程序结束时关闭日志文件
		defer fileLogger.Close()
	}

	if global.GVA_DB != nil {
		core.MigrateModels()
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}

	app := &cli.App{
		Name:  "GDocs",
		Usage: "Documentation Service",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Start the docs server",
				Action: func(c *cli.Context) error {
					core.RunServer()
					return nil
				},
			},
			{
				Name:    "createsuperuser",
				Aliases: []string{"csu"},
				Usage:   "Create a superuser",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username",
						Aliases:  []string{"u"},
						Usage:    "Superuser username",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Usage:    "Superuser password",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					username := c.String("username")
					password := c.String("password")
					return commands.CreateSuperUser(username, password)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		global.GVA_LOG.Error("commands start err", err)
	}
}
