package admin

import (
	"fmt"
	"net/http"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"github.com/bookandmusic/docs/common"
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
	site_info := common.GenerateSiteInfo()
	c.HTML(http.StatusOK, "admin/index.html", pongo2.Context{
		"serverName":  site_info.Name,
		"keywords":    site_info.Keywords,
		"description": site_info.Desc,
		"site_icon":   site_info.Icon,
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
	site_info := common.GenerateSiteInfo()
	// 创建 MenuData 结构体实例，表示菜单数据
	menu := MenuItemData{
		HomeInfo: MenuItem{
			Title: "首页",
			Href:  "welcome/",
		},
		LogoInfo: MenuItem{
			Title: site_info.Name,
			Image: site_info.Logo,
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
					{
						Title:  "系统设置",
						Href:   "system/settings/",
						Icon:   "fa fa-cog",
						Target: "_self",
					},
				},
			},
		},
	}

	// 将菜单数据转换为 JSON 格式并返回给客户端
	c.JSON(http.StatusOK, menu)
}

func (u *AdminController) SystemSetting(c *gin.Context) {
	if c.Request.Method == "GET" {
		_, giscus_info := common.GenerateGiscusInfo()
		c.HTML(http.StatusOK, "admin/system/settings.html", pongo2.Context{
			"csrfToken":   csrf.Token(c.Request),
			"site_info":   common.GenerateSiteInfo(),
			"person_info": common.GeneratePersonInfo(),
			"giscus_info": giscus_info,
		})
	} else {
		resource := c.Query("resource")
		var (
			obj     map[string]string
			err_msg string
		)
		if err := c.BindJSON(&obj); err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("parse params err: %v", err))
			switch resource {
			case "site_info":
				err_msg = "站点信息有误"
			case "person_info":
				err_msg = "个人信息有误"
			case "giscus_info":
				err_msg = "评论Giscus插件信息有误"
			}
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": err_msg})
			return
		}

		// 遍历结构体字段
		for key, val := range obj {
			if key != "" && val != "" {
				if err := models.NewSetting().UpdateOrCreate(key, val); err != nil {
					global.GVA_LOG.Error(fmt.Sprintf("update system setting '%s':'%s', err: %v", key, val, err))
				}
			} else if key != "" && val == "" {
				if err := models.NewSetting().DeleteByName(key); err != nil {
					global.GVA_LOG.Error(fmt.Sprintf("update system setting '%s' default val, err: %v", key, err))
				}
			}
		}
	}
}
