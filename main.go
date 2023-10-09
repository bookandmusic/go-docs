package main

import (
	"os"
	"strconv"

	cli "github.com/urfave/cli/v2"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/bookandmusic/docs/commands"
	"github.com/bookandmusic/docs/core"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

var (
	configFile string
	serverPort int = 8080

	username string
	password string
)

func serverInitAction() {
	// 初始化配置
	global.GVA_VP = core.InitConfig(configFile, serverPort)
	global.GVA_LOG = core.InitLog()
	// 初始化数据库连接
	global.GVA_DB = core.InitDatabase()
	global.GVA_MINIFY = core.InitMinify()
	global.GVA_BLEVE_INDEX = core.InitBleveIndex()

	if global.GVA_DB != nil {
		core.MigrateModels()
	}

	if global.GVA_BLEVE_INDEX != nil && global.GVA_DB != nil {
		models.NewArticle().InitArticleIndex()
	}
}

func dbInitAction() {
	// 初始化配置
	global.GVA_VP = core.InitConfig(configFile, serverPort)
	global.GVA_LOG = core.InitLog()
	// 初始化数据库连接
	global.GVA_DB = core.InitDatabase()
}

func main() {
	app := &cli.App{
		Name:  "GDocs",
		Usage: "Documentation Service",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Start the docs server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "",
						Usage:   "Load configuration from `FILE`",
					},
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "",
						Usage:   "Server port",
					},
				},
				Action: func(c *cli.Context) error {
					configFile = c.String("config")
					if port, err := strconv.Atoi(c.String("port")); err == nil {
						if port >= 5000 && port <= 65535 {
							serverPort = port
						}
					}

					serverInitAction()
					// 程序结束前关闭数据库链接
					db, _ := global.GVA_DB.DB()
					defer db.Close()
					if fileLogger, ok := global.GVA_LOG.Out.(*lumberjack.Logger); ok {
						// 在程序结束时关闭日志文件
						defer fileLogger.Close()
					}
					commands.CreateSuperUser("admin", "123456")
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
					dbInitAction()
					// 程序结束前关闭数据库链接
					db, _ := global.GVA_DB.DB()
					defer db.Close()
					if fileLogger, ok := global.GVA_LOG.Out.(*lumberjack.Logger); ok {
						// 在程序结束时关闭日志文件
						defer fileLogger.Close()
					}
					return commands.CreateSuperUser(username, password)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		global.GVA_LOG.Error("commands start err", err)
	}
}
