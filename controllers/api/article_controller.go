package api

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
	"github.com/bookandmusic/docs/utils"
)

type ArticleAPIController struct{}

type articleJson struct {
	ID           int      `json:"id"`
	Title        string   `json:"title" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	CollectionId int      `json:"collectionId"`
	CategoryId   int      `json:"categoryId"`
	Tags         []string `json:"tags"`
}

func NewArticleAPIController() *ArticleAPIController {
	return &ArticleAPIController{}
}

func (controller *ArticleAPIController) ArticleDetail(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	var (
		article *models.Article
		err     error
	)

	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil || articleId <= 0 {
		c.JSON(http.StatusOK, common.ParamError)
		return
	}

	if article, err = models.NewArticle().FindByArticleId(articleId); err != nil || article == nil {
		global.GVA_LOG.Warn(fmt.Sprintf("find article by id err: %v", err))
		c.JSON(http.StatusOK, common.NotExistsError)
		return
	}
	var tags []string
	for _, tag := range article.Tags {
		tags = append(tags, tag.Name)

	}
	data := map[string]any{
		"id":           article.ID,
		"title":        article.Title,
		"content":      article.Content,
		"type":         article.Type.String(),
		"tags":         tags,
		"categoryId":   article.CategoryID,
		"collectionId": article.CollectionID,
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: data})
}

func (controller *ArticleAPIController) ArticleList(c *gin.Context) {
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	typeStr := c.Query("type")
	pageStr := c.Query("page")
	sizeStr := c.Query("limit")
	categoryIdStr := c.Query("categoryId")
	collectionIdStr := c.Query("collectionId")

	// 默认值
	defaultPage := 1
	defaultSize := 10
	var articleType models.ArticleTyle = -1

	// 将查询参数转换为整数，如果无法转换或不存在，则使用默认值
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = defaultPage
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = defaultSize
	}
	if typeStr != "" {
		articleType = models.ArticleTyle.Parse(models.Blog, typeStr)
	}

	categoryID, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		categoryID = -1
	}

	collectionID, err := strconv.Atoi(collectionIdStr)
	if err != nil || (articleType != models.Doc) {
		collectionID = -1
	}

	var (
		articles []*models.Article
		total    int
	)
	params := map[string]any{
		"page":         page,
		"size":         size,
		"keyword":      keyword,
		"sort":         sort,
		"articleType":  articleType,
		"categoryID":   categoryID,
		"collectionID": collectionID,
	}
	articles, total, _ = models.NewArticle().FindArticlesByParams(params)
	var data []map[string]any
	for _, article := range articles {
		var tagNames []string
		for _, tag := range article.Tags {
			tagNames = append(tagNames, tag.Name)
		}
		data = append(data, map[string]any{
			"ID":             article.ID,
			"UpdatedAt":      article.UpdatedAt,
			"title":          article.Title,
			"type":           article.Type.String(),
			"categoryName":   article.Category.Name,
			"categoryId":     article.Category.ID,
			"collectionName": article.Collection.Name,
			"collectionId":   article.Collection.ID,
			"tags":           tagNames,
			"tagsNum":        len(tagNames),
		})
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: map[string]any{"total": total, "data": data}})
}

func (controller *ArticleAPIController) DeleteArticle(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil || articleId <= 0 {
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	var articleIDs []int
	articleIDs = append(articleIDs, articleId)

	if err := models.NewArticle().DeleteByArticleIds(articleIDs); err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("delete article err: %v", err))
		c.JSON(http.StatusOK, common.DeleteError)
		return
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: ""})
}

func (controller *ArticleAPIController) EditArticle(c *gin.Context) {
	// 从路由中获取整数参数
	var (
		obj  *models.Article
		json articleJson
	)

	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		global.GVA_LOG.Error(fmt.Sprintf("参数解析失败: %v", err))
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	articleType := models.ArticleTyle.Parse(models.Blog, json.Type)
	switch articleType {
	case models.Blog:
		if json.CollectionId != 0 {
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "博客不能属于文集"})
			return
		}
		if json.CategoryId == 0 {
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "博客必须属于一个分类"})
			return
		}
	case models.Doc:
		if json.CategoryId == 0 {
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "文档必须属于一个文集"})
			return
		}
	default:
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "文档类别不对"})
		return

	}

	// 将查询参数转换为整数，如果无法转换或不存在，则使用默认值
	var (
		collectionID int      = json.CollectionId
		categoryID   int      = json.CategoryId
		tagsName     []string = json.Tags
		title        string   = json.Title
		err          error

		collection *models.Collection = models.NewCollection()
		category   *models.Category   = models.NewCategory()
		tags       []models.Tag
		tag        *models.Tag

		htmlContent string
		frontMatter *models.FrontMatter
		markdownStr string
	)

	if collectionID != 0 {
		if collection, err = models.NewCollection().FindByCollectionId(collectionID); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("查询文集失败: %v", err))
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "文集不存在"})
			return
		}
	}

	if categoryID != 0 {
		if category, err = models.NewCategory().FindByCategoryId(categoryID); err != nil {
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "分类不存在"})
			return
		}
	}
	for _, tagName := range tagsName {
		tag, err = models.NewTag().FindOrCreateByName(tagName)
		if err != nil {
			continue
		}
		tags = append(tags, *tag)
	}

	if htmlContent, markdownStr, frontMatter, err = models.NewArticle().MarkdownToHTML(json.Content); err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("转换markdown失败: %v", err))
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "添加文档失败"})
		return
	}

	if frontMatter != nil {
		title = frontMatter.Title
	}

	if json.ID == 0 {
		obj = &models.Article{
			Title:       title,
			Content:     markdownStr,
			HtmlContent: htmlContent,
			Type:        articleType,
			Category:    *category,
			Collection:  *collection,
			Tags:        tags,
		}
		if err := models.NewArticle().Create(obj); err != nil {
			global.GVA_LOG.Warn("Failed to add article", err)
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "添加文档失败"})
			return
		}
	} else {
		if obj, err = models.NewArticle().FindByArticleId(json.ID); err != nil {
			global.GVA_LOG.Warn("Failed to find article", err)
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "查询文档失败"})
			return
		}
		updates := map[string]interface{}{
			"title":        title,
			"content":      markdownStr,
			"html_content": htmlContent,
			"type":         articleType,
			"Category":     category,
			"Tags":         tags,
			"Collection":   collection,
		}
		if err := models.NewArticle().Update(obj, updates); err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to edit article: %v", err))
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "更新文档失败"})
			return
		}
	}

	var tagNames []string
	for _, tag := range obj.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	var data = map[string]any{
		"ID":             obj.ID,
		"UpdatedAt":      obj.UpdatedAt,
		"title":          obj.Title,
		"type":           obj.Type.String(),
		"categoryName":   obj.Category.Name,
		"categoryId":     obj.Category.ID,
		"collectionName": obj.Collection.Name,
		"collectionId":   obj.Collection.ID,
		"tags":           tagNames,
		"tagsNum":        len(tagNames),
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: data})
}

func (controller *ArticleAPIController) ImportMds(c *gin.Context) {

	// 获取文集 ID
	collectionIdStr := c.PostForm("collectionId")
	articleTypeStr := c.PostForm("articleType")
	categoryIdStr := c.PostForm("categoryId")
	var (
		collectionID, categoryID int
		articleType              models.ArticleTyle
		category                 *models.Category
		collection               *models.Collection
		err                      error
	)
	articleType = models.Blog.Parse(articleTypeStr)

	if articleTypeStr == "" {
		articleTypeStr = "blog"
	}
	if articleType == models.Blog {
		if categoryIdStr == "" {
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "博客必须属于一个分类"})
			return
		}
		if categoryID, err = strconv.Atoi(categoryIdStr); err != nil {
			c.JSON(http.StatusOK, common.ParamError)
			return
		}
		category, err = models.NewCategory().FindByCategoryId(categoryID)
		if category == nil || err != nil {
			c.JSON(http.StatusOK, common.ParamError)
			return
		}
	} else {
		if collectionIdStr == "" {
			c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "文档必须属于一个文集"})
			return
		}
		if collectionID, err = strconv.Atoi(collectionIdStr); err != nil {
			c.JSON(http.StatusOK, common.ParamError)
			return
		}
		collection, err = models.NewCollection().FindByCollectionId(collectionID)
		if collection == nil || err != nil {
			c.JSON(http.StatusOK, common.ParamError)
			return
		}
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "文件上传失败"})
		return
	}

	// 检查文件类型是否为 ZIP 文件
	if filepath.Ext(file.Filename) != ".zip" {
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ParamErrorCode, Message: "文件类型错误"})
		return
	}

	// 创建基于会话 ID 的临时目录
	SiteName := models.NewSetting().GetValue("site_name")
	if SiteName == "" {
		SiteName = "GDocs"
	}
	timeStr := utils.GenerateMD5Hash(time.Now().String())
	tmpDir := path.Join(os.TempDir(), SiteName, timeStr)
	//  确保临时目录存在
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "无法创建临时目录"})
		return
	}
	// 保存上传的文件到临时目录
	filePath := path.Join(tmpDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "无法保存文件"})
		return
	}

	// 导入zip项目
	if err := models.NewArticle().ImportMds(filePath, articleType, collectionID, categoryID); err != nil {
		global.GVA_LOG.Error("upload markdown zip error", err)
		c.JSON(http.StatusOK, common.ErrorMsg{Code: common.ServerErrorCode, Message: "文件处理失败"})
	}

	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: ""})
}
