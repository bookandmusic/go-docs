package common

import (
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type GiscusConfig struct {
	GiscusDataRepo       string `json:"giscus_repo"`
	GiscusDataRepoId     string `json:"giscus_repo_id"`
	GiscusDataCategory   string `json:"giscus_category"`
	GiscusDataCategoryId string `json:"giscus_category_id"`
}

type SiteInfo struct {
	Name     string `json:"site_name"`
	Url      string `json:"site_url"`
	Desc     string `json:"site_description"`
	Keywords string `json:"site_keyword"`
	Icon     string `json:"site_icon"`
	Logo     string `json:"site_logo"`
	Since    string `json:"since_year"`
	Beian    string `json:"site_beian"`
	Comment  string `json:"comment"`
}

type PersonInfo struct {
	AuthorName    string `json:"author_name"`
	AuthorAddress string `json:"author_address"`
	AuthorAvatar  string `json:"author_avatar"`
	Email         string `json:"email"`
	Github        string `json:"github"`
	Wechat        string `json:"weichat"`
	Weibo         string `json:"weibo"`
	QQ            string `json:"qq"`
}

type ArticleInfo struct {
	ArticleCount    int `json:"article_count"`
	TagCount        int `json:"tag_count"`
	CategoryCount   int `json:"category_count"`
	CollectionCount int `json:"collection_count"`
}

func GenerateSiteInfo() SiteInfo {
	SinceYear := models.NewSetting().GetValue("since_year")
	SiteUrl := models.NewSetting().GetValue("site_url")
	SiteIcon := models.NewSetting().GetValue("site_icon")
	SiteLogo := models.NewSetting().GetValue("site_logo")
	SiteName := models.NewSetting().GetValue("site_name")
	SiteDescription := models.NewSetting().GetValue("site_description")
	SiteKeyword := models.NewSetting().GetValue("site_keyword")
	SiteBeian := models.NewSetting().GetValue("site_beian")
	Comment := models.NewSetting().GetValue("comment")

	if SiteIcon == "" {
		SiteIcon = "/static/public/images/favicon.ico"
	}

	if SiteLogo == "" {
		SiteLogo = "/static/public/images/logo.png"
	}

	if SiteName == "" {
		SiteName = global.GVA_CONFIG.Server.ServerName
	}

	if SiteDescription == "" {
		SiteDescription = global.GVA_CONFIG.Server.Description
	}

	if SiteKeyword == "" {
		SiteKeyword = global.GVA_CONFIG.Server.Keyword
	}

	site_info := SiteInfo{
		Name:     SiteName,
		Desc:     SiteDescription,
		Keywords: SiteKeyword,
		Since:    SinceYear,
		Icon:     SiteIcon,
		Logo:     SiteLogo,
		Url:      SiteUrl,
		Beian:    SiteBeian,
		Comment:  Comment,
	}
	return site_info
}

func GenerateGiscusInfo() (bool, GiscusConfig) {
	GiscusDataRepo := models.NewSetting().GetValue("giscus_repo")
	GiscusDataRepoId := models.NewSetting().GetValue("giscus_repo_id")
	GiscusDataCategory := models.NewSetting().GetValue("giscus_category")
	GiscusDataCategoryId := models.NewSetting().GetValue("giscus_category_id")
	var isComment bool
	if GiscusDataCategory != "" && GiscusDataCategoryId != "" && GiscusDataRepo != "" && GiscusDataRepoId != "" {
		isComment = true
	} else {
		isComment = false
	}
	return isComment, GiscusConfig{
		GiscusDataRepo:       GiscusDataRepo,
		GiscusDataRepoId:     GiscusDataRepoId,
		GiscusDataCategory:   GiscusDataCategory,
		GiscusDataCategoryId: GiscusDataCategoryId,
	}
}

func GenerateArticleInfo() ArticleInfo {
	cateTotal := models.NewCategory().Count()
	tagTotal := models.NewTag().Count()
	articleTotal := models.NewArticle().Count()
	collectionTotal := models.NewCollection().Count()

	return ArticleInfo{
		ArticleCount:    articleTotal,
		TagCount:        tagTotal,
		CategoryCount:   cateTotal,
		CollectionCount: collectionTotal,
	}
}

func GeneratePersonInfo() PersonInfo {
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
		AuthorName:    authorName,
		AuthorAddress: authorAddress,
		AuthorAvatar:  authorAvatar,
		Email:         email,
		Github:        github,
		Wechat:        weichat,
		QQ:            qq,
		Weibo:         weibo,
	}
	return person_info
}
