package admin

import (
	"net/http"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

type MenuItem struct {
	Title  string     `json:"title"`
	Icon   string     `json:"icon"`
	Image  string     `json:"image"`
	Href   string     `json:"href"`
	Target string     `json:"target"`
	Child  []MenuItem `json:"child"`
}

type MenuItemData struct {
	HomeInfo MenuItem   `json:"homeInfo"`
	LogoInfo MenuItem   `json:"logoInfo"`
	MenuInfo []MenuItem `json:"menuInfo"`
}

// 定义一个结构体来表示菜单数据

func (u *AdminController) Index(c *gin.Context) {
	SiteIcon := models.NewSetting().GetValue("site_icon")

	if SiteIcon == "" {
		SiteIcon = "/static/public/images/favicon.ico"
	}

	c.HTML(http.StatusOK, "admin/index.html", pongo2.Context{
		"serverName":  global.GVA_CONFIG.Server.ServerName,
		"keywords":    global.GVA_CONFIG.Server.KeyWord,
		"description": global.GVA_CONFIG.Server.Description,
		"site_icon":   SiteIcon,
	})
}

func (u *AdminController) Welcome(c *gin.Context) {
	cateTotal := models.NewCategory().Count()
	tagTotal := models.NewTag().Count()
	articleTotal := models.NewArticle().Count()
	collectionTotal := models.NewCollection().Count()
	journals, _ := models.NewJournal().FindAll()
	if len(journals) > 5 {
		journals = journals[:5]
	}
	c.HTML(http.StatusOK, "admin/welcome.html", pongo2.Context{
		"cateTotal":       cateTotal,
		"tagTotal":        tagTotal,
		"articleTotal":    articleTotal,
		"collectionTotal": collectionTotal,
		"journals":        journals,
	})
}

func (u *AdminController) AdminCenterMenu(c *gin.Context) {
	SiteLogo := models.NewSetting().GetValue("site_logo")

	if SiteLogo == "" {
		SiteLogo = "/static/public/images/logo.png"
	}

	// 创建 MenuData 结构体实例，表示菜单数据
	menu := MenuItemData{
		HomeInfo: MenuItem{
			Title: "首页",
			Href:  "welcome/",
		},
		LogoInfo: MenuItem{
			Title: global.GVA_CONFIG.Server.ServerName,
			Image: SiteLogo,
			Href:  "",
		},
		MenuInfo: []MenuItem{
			{
				Title:  "仪表盘",
				Href:   "welcome/",
				Icon:   "fa fa-dashboard",
				Target: "_self",
			},
			{
				Title:  "后台管理",
				Icon:   "fa fa-navicon",
				Href:   "",
				Target: "_self",
				Child: []MenuItem{
					{
						Title:  "分类",
						Href:   "category/",
						Icon:   "fa fa-folder-open",
						Target: "_self",
					},
					{
						Title:  "标签",
						Href:   "tag/",
						Icon:   "fa fa-tags",
						Target: "_self",
					},
					{
						Title:  "文集",
						Href:   "collection/",
						Icon:   "fa fa-book",
						Target: "_self",
					},
					{
						Title:  "文章",
						Href:   "article/",
						Icon:   "fa fa-newspaper-o",
						Target: "_self",
					},
					{
						Title:  "日志",
						Href:   "journal/",
						Icon:   "fa fa-pencil-square",
						Target: "_self",
					},
				},
			},
			{
				Title:  "系统工具",
				Icon:   "fa fa-cogs",
				Href:   "",
				Target: "_self",
				Child: []MenuItem{
					{
						Title:  "文档导入",
						Href:   "system/import_mds/",
						Icon:   "fa fa-hand-o-right",
						Target: "_self",
					},
				},
			},
		},
	}

	// 将菜单数据转换为 JSON 格式并返回给客户端
	c.JSON(http.StatusOK, menu)
}
