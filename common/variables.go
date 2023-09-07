package common

import (
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type SiteInfo struct {
	Name     string
	Url      string
	Desc     string
	Keywords string
	Icon     string
	Logo     string
	Since    string
}

type PersonInfo struct {
	AuthorName      string
	AuthorAddress   string
	AuthorAvatar    string
	Email           string
	Github          string
	Wechat          string
	Weibo           string
	QQ              string
	ArticleCount    int
	TagCount        int
	CategoryCount   int
	CollectionCount int
}

func GenerateSiteInfo() SiteInfo {
	setting := global.GVA_CONFIG.Server

	SinceYear := models.NewSetting().GetValue("since_year")
	SiteUrl := models.NewSetting().GetValue("site_url")
	SiteIcon := models.NewSetting().GetValue("site_icon")
	SiteLogo := models.NewSetting().GetValue("site_logo")

	if SiteIcon == "" {
		SiteIcon = "/static/public/images/favicon.ico"
	}

	if SiteLogo == "" {
		SiteLogo = "/static/public/images/logo.png"
	}

	site_info := SiteInfo{
		Name:     setting.ServerName,
		Desc:     setting.Description,
		Keywords: setting.KeyWord,
		Since:    SinceYear,
		Icon:     SiteIcon,
		Logo:     SiteLogo,
		Url:      SiteUrl,
	}
	return site_info
}

func GeneratePersonInfo() PersonInfo {
	cateTotal := models.NewCategory().Count()
	tagTotal := models.NewTag().Count()
	articleTotal := models.NewArticle().Count()
	collectionTotal := models.NewCollection().Count()
	authorName := models.NewSetting().GetValue("author_name")
	authorAddress := models.NewSetting().GetValue("author_address")
	authorAvatar := models.NewSetting().GetValue("author_avatar")
	email := models.NewSetting().GetValue("email")
	qq := models.NewSetting().GetValue("qq")
	weichat := models.NewSetting().GetValue("weichat")
	weibo := models.NewSetting().GetValue("weibo")
	github := models.NewSetting().GetValue("github")

	if authorName == "" {
		authorName = "Mr.Liu"
	}
	if authorAddress == "" {
		authorAddress = "北京/朝阳"
	}

	if authorAvatar == "" {
		authorAvatar = "/static/public/images/avatar.png"
	}

	person_info := PersonInfo{
		AuthorName:      authorName,
		AuthorAddress:   authorAddress,
		AuthorAvatar:    authorAvatar,
		Email:           email,
		Github:          github,
		Wechat:          weichat,
		QQ:              qq,
		Weibo:           weibo,
		ArticleCount:    articleTotal,
		TagCount:        tagTotal,
		CategoryCount:   cateTotal,
		CollectionCount: collectionTotal,
	}
	return person_info
}
