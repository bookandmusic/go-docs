package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
	"github.com/bookandmusic/docs/utils"
)

type ArticleController struct{}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

func (controller *ArticleController) Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/article/admin.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
		})
	}
}

func (controller *ArticleController) ArticleList(c *gin.Context) {
	keyword := c.Query("keyword")
	typeStr := c.Query("type")
	pageStr := c.Query("page")
	sizeStr := c.Query("limit")
	categoryIdStr := c.Query("categoryId")
	collectionIdStr := c.Query("collectionId")

	// 默认值
	defaultPage := 1
	defaultSize := 10

	// 将查询参数转换为整数，如果无法转换或不存在，则使用默认值
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = defaultPage
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = defaultSize
	}
	articleType, err := strconv.Atoi(typeStr)
	if err != nil || (articleType != 0 && articleType != 1) {
		articleType = -1
	}

	categoryID, err := strconv.Atoi(categoryIdStr)
	if err != nil || (articleType != 0 && articleType != 1) {
		categoryID = -1
	}

	collectionID, err := strconv.Atoi(collectionIdStr)
	if err != nil || (articleType != 0 && articleType != 1) {
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
		"articleType":  articleType,
		"categoryID":   categoryID,
		"collectionID": collectionID,
	}
	articles, total, _ = models.NewArticle().FindArticlesByParams(params)
	var data []map[string]any
	for _, article := range articles {
		data = append(data, map[string]any{
			"id":              article.ID,
			"title":           article.Title,
			"category_name":   article.Category.Name,
			"collection_name": article.Collection.Name,
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": total, "data": data})
	return
}

func (controller *ArticleController) DeleteArticle(c *gin.Context) {
	type paramsStruct struct {
		ArticleIDs *[]int `json:"articleIDs"`
	}
	var params paramsStruct
	if err := c.BindJSON(&params); err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("parse params err: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文章ID无效"})
		return
	}

	err := models.NewArticle().DeleteByArticleIds(*params.ArticleIDs)
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("delete article err: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除文档失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除文档成功"})
	return
}

func (controller *ArticleController) EditArticle(c *gin.Context) {
	// 从路由中获取整数参数
	var obj *models.Article
	articleIDStr := c.Query("id")
	articleTypeStr := c.Query("article_type")
	if articleIDStr == "" || articleIDStr == "0" {
		obj = models.NewArticle()
	} else {
		// 将字符串形式的 articleID 转换为整数
		articleID, err := strconv.Atoi(articleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "无效的文档ID"})
			return
		}
		obj, _ = models.NewArticle().FindByArticleId(articleID)
	}

	if articleTypeStr == "" {
		articleTypeStr = "blog"
	}

	if c.Request.Method == "GET" {
		// 处理 GET 请求，返回登录页面
		tagIds := []int{}
		for _, tag := range obj.Tags {
			tagIds = append(tagIds, int(tag.ID))
		}
		jsonData, _ := json.Marshal(tagIds)
		var template_name string
		if articleTypeStr == "blog" {
			template_name = "admin/article/blog_edit.html"
		} else {
			template_name = "admin/article/doc_edit.html"
		}
		c.HTML(http.StatusOK, template_name, pongo2.Context{
			"csrfToken":      csrf.Token(c.Request),
			"article":        obj,
			"tagIDs":         string(jsonData),
			"articleTypeStr": articleTypeStr,
			"articleContent": utils.EscapeJS(obj.Content),
		})
	} else if c.Request.Method == "POST" {
		title := c.PostForm("title")
		content := c.PostForm("content")
		articleTypeStr := c.PostForm("type")
		collectionIdStr := c.PostForm("collectionId")
		categoryIdStr := c.PostForm("categoryId")
		tagIDsStr := c.PostForm("tagIds")

		if title == "" || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "缺少必要内容"})
			return
		}

		if articleTypeStr != "0" && articleTypeStr != "1" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文档类别不对"})
			return
		}
		articleType, _ := strconv.Atoi(articleTypeStr)

		if articleType == 1 && collectionIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文档必须属于一个文集"})
			return
		}

		if articleType == 0 && collectionIdStr != "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "博客不能属于文集"})
			return
		}
		if articleType == 0 && categoryIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "博客分类不能为空"})
			return
		}

		// 将查询参数转换为整数，如果无法转换或不存在，则使用默认值
		var (
			collectionID int
			categoryID   int
			tagID        int
			tagIDs       []int
			collection   *models.Collection = models.NewCollection()
			category     *models.Category   = models.NewCategory()
			tags         []models.Tag
			err          error

			htmlContent string
			tocJson     string
			frontMatter *models.FrontMatter
			markdownStr string
		)
		if collectionIdStr != "" {
			if collectionID, err = strconv.Atoi(collectionIdStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文集ID无效"})
				return
			}
		}

		if categoryIdStr != "" {
			if categoryID, err = strconv.Atoi(categoryIdStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "分类ID无效"})
				return
			}
		}

		// 使用逗号分割字符串
		if tagIDsStr != "" {
			strTagIds := strings.Split(tagIDsStr, ",")
			for _, idStr := range strTagIds {
				if tagID, err = strconv.Atoi(idStr); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "标签ID无效"})
					return
				}
				tagIDs = append(tagIDs, tagID)
			}
		}

		if collectionID != 0 {
			if collection, err = models.NewCollection().FindByCollectionId(collectionID); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("查询文集失败: %v", err))
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文集不存在"})
				return
			}
		}

		if categoryID != 0 {
			if category, err = models.NewCategory().FindByCategoryId(categoryID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "分类不存在"})
				return
			}
		}
		if len(tagIDs) > 0 {
			if tags, err = models.NewTag().FindByTagIdArray(tagIDs); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "标签不存在"})
				return
			}
		}

		if htmlContent, markdownStr, frontMatter, err = models.NewArticle().MarkdownToHTML(content); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("转换markdown失败: %v", err))
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "添加文档失败"})
			return
		}

		if frontMatter != nil {
			title = frontMatter.Title
		}

		if obj == nil || obj.ID == 0 {
			article := &models.Article{
				Title:       title,
				Content:     markdownStr,
				HtmlContent: htmlContent,
				Type:        models.ArticleTyle(articleType),
				Category:    *category,
				Collection:  *collection,
				Tags:        tags,
			}
			if err := models.NewArticle().Create(article); err != nil {
				global.GVA_LOG.Warn("Failed to add article", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "添加文档失败"})
				return
			}
		} else {
			updates := map[string]interface{}{
				"title":        title,
				"content":      markdownStr,
				"html_content": htmlContent,
				"toc":          tocJson,
				"type":         models.ArticleTyle(articleType),
				"Category":     category,
				"Tags":         tags,
				"Collection":   collection,
			}
			if err := obj.Update(updates); err != nil {
				global.GVA_LOG.Warn("Failed to edit article", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "更新文档失败"})
				return
			}
		}

		return
	}
}

func (controller *ArticleController) ImportMds(c *gin.Context) {
	// 处理 GET 请求
	if c.Request.Method == http.MethodGet {
		// 查询所有文集
		collections, err := models.NewCollection().FindByKeyword("", false, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "无法获取文集"})
			return
		}

		categories, err := models.NewCategory().FindByKeyword("")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "无法获取文集"})
			return
		}

		// 渲染 HTML 模板
		c.HTML(http.StatusOK, "admin/article/import_mds.html", pongo2.Context{
			"csrfToken":   csrf.Token(c.Request),
			"collections": collections,
			"categories":  categories,
		})
		return
	}

	// 处理 POST 请求
	if c.Request.Method == http.MethodPost {
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
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "博客分类不能为空"})
				return
			}
			if categoryID, err = strconv.Atoi(categoryIdStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "分类ID无效"})
				return
			}
			category, err = models.NewCategory().FindByCategoryId(categoryID)
			if category == nil || err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "分类ID无效"})
				return
			}
		} else {
			if collectionIdStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文档必须属于一个文集"})
				return
			}
			if collectionID, err = strconv.Atoi(collectionIdStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文集ID无效"})
				return
			}
			collection, err = models.NewCollection().FindByCollectionId(collectionID)
			if collection == nil || err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文集ID无效"})
				return
			}
		}

		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文件上传失败"})
			return
		}

		// 检查文件类型是否为 ZIP 文件
		if filepath.Ext(file.Filename) != ".zip" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "文件类型错误"})
			return
		}

		// 创建基于会话 ID 的临时目录
		timeStr := utils.GenerateMD5Hash(time.Now().String())
		tmpDir := path.Join(os.TempDir(), timeStr)
		//  确保临时目录存在
		if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "无法创建临时目录"})
			return
		}
		// 保存上传的文件到临时目录
		filePath := path.Join(tmpDir, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "无法保存文件"})
			return
		}

		// 导入zip项目
		if err := models.NewArticle().ImportMds(filePath, articleType, collectionID, categoryID); err != nil {
			global.GVA_LOG.Error("upload markdown zip error", err)
			c.JSON(http.StatusOK, gin.H{"status": true, "msg": "文件处理失败"})
		}

		// 如果一切顺利，响应成功消息
		c.JSON(http.StatusOK, gin.H{"status": true, "msg": "文件上传并处理成功"})
	}
}
