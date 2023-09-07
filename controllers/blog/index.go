package blog

import (
	"fmt"
	"net/http"
	"strconv"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
	"github.com/bookandmusic/docs/utils"
)

type BlogController struct{}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (controller *BlogController) Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		pageStr := c.Param("page")
		categoryStr := c.Param("category")
		tagStr := c.Param("tag")
		var (
			err      error
			page     int
			category *models.Category = models.NewCategory()
			tag      *models.Tag      = models.NewTag()
			pre      bool
			next     bool
			title    string
			articles []*models.Article
			total    int
			pageType string
		)
		if page, err = strconv.Atoi(pageStr); err != nil {
			page = 1
		}
		site_info := common.GenerateSiteInfo()
		person_info := common.GeneratePersonInfo()

		if categoryStr != "" && tagStr == "" {
			category, err = models.NewCategory().FindByIdentify(categoryStr)
			title = category.Name
			if articles, total, err = models.NewArticle().FindBlogsByCategoryID(page, 10, int(category.ID)); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("find blogs by categoryID error: %v", err))
			}
			pageType = "分类"
		} else if categoryStr == "" && tagStr != "" {
			tag, _ = models.NewTag().FindByIdentify(tagStr)
			title = tag.Name
			if articles, total, err = models.NewArticle().FindBlogsByTagID(page, 10, int(tag.ID)); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("find blogs by tagID error: %v", err))
			}
			pageType = "标签"
		} else {
			title = "首页"
			if articles, total, err = models.NewArticle().FindBlogs(page, 10); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("find blogs error: %v", err))
			}
		}

		tags, _ := models.NewTag().FindAll()

		// 计算总页数
		totalPages := int((total + 10 - 1) / 10)
		pages := utils.GeneratePageNumbers(totalPages, page)

		if page == 1 {
			pre = false
		} else {
			pre = true
		}
		if page == totalPages || totalPages == 0 {
			next = false
		} else {
			next = true
		}

		c.HTML(http.StatusOK, "blog/index.html", pongo2.Context{
			"site_info":    site_info,
			"page_title":   title,
			"person_info":  person_info,
			"tags":         tags,
			"articles":     articles,
			"page_type":    pageType,
			"pages":        pages,
			"current_page": page,
			"pre":          pre,
			"next":         next,
		})
	}
}

func (controller *BlogController) Post(c *gin.Context) {
	if c.Request.Method == "GET" {
		identify := c.Param("identify")
		obj := models.NewArticle()
		article, _ := obj.FindByIdentify(identify)
		site_info := common.GenerateSiteInfo()
		person_info := common.GeneratePersonInfo()
		if article == nil {
			c.HTML(http.StatusNotFound, "public/404.html", pongo2.Context{
				"site_info": site_info,
				"err_msg":   "查询不到该文章的信息", // 获取当前路由地址
			})
		} else {
			if article.CollectionID == 0 {
				previousArticle, nextArticle := article.FindPreviousBlogAndNextBlog(int(article.ID))
				c.HTML(http.StatusOK, "blog/post.html", pongo2.Context{
					"article":          article,
					"site_info":        site_info,
					"page_title":       article.Title,
					"person_info":      person_info,
					"previous_article": previousArticle,
					"next_article":     nextArticle,
				})
			} else {
				articles, _ := obj.FindByCollectionId(int(article.CollectionID))
				tocList := obj.TocList(articles)
				previousArticle, nextArticle := article.FindPreviousDocAndNextDoc(tocList, int(article.ID))
				c.HTML(http.StatusOK, "blog/article.html", pongo2.Context{
					"article":          article,
					"toc_list":         tocList,
					"site_info":        site_info,
					"page_title":       article.Title,
					"person_info":      person_info,
					"previous_article": previousArticle,
					"next_article":     nextArticle,
				})
			}

		}
	}
}

func (controller *BlogController) Collections(c *gin.Context) {
	if c.Request.Method == "GET" {
		site_info := common.GenerateSiteInfo()
		title := "文集"

		collections, _ := models.NewCollection().FindByKeyword("", false, true)

		c.HTML(http.StatusOK, "blog/collections.html", pongo2.Context{
			"site_info":   site_info,
			"page_title":  title,
			"collections": collections,
		})
	}
}

func (controller *BlogController) Archives(c *gin.Context) {
	if c.Request.Method == "GET" {
		site_info := common.GenerateSiteInfo()
		person_info := common.GeneratePersonInfo()
		title := "归档"

		articles, _ := models.NewArticle().FindBlogItems()

		// 使用 map 进行归档
		title_articles := make(map[string][]*models.ArticleItem)
		for _, article := range articles {
			title := article.CreatedAt.Format("2006年01月")
			title_articles[title] = append(title_articles[title], article)
		}

		c.HTML(http.StatusOK, "blog/timeline.html", pongo2.Context{
			"site_info":      site_info,
			"page_title":     title,
			"person_info":    person_info,
			"title_articles": title_articles,
		})
	}
}

func (controller *BlogController) Categories(c *gin.Context) {
	if c.Request.Method == "GET" {
		site_info := common.GenerateSiteInfo()
		person_info := common.GeneratePersonInfo()
		title := "分类"

		categories, _ := models.NewCategory().FindAll()

		categories_map := make(map[int]models.Category)
		for _, category := range categories {
			categories_map[int(category.ID)] = *category
		}
		articles, _ := models.NewArticle().FindBlogItems()

		// 使用 map 进行归档
		title_articles := make(map[string][]*models.ArticleItem)
		for _, article := range articles {
			cateID := article.CategoryID
			var cate string
			if cateID == 0 {
				cate = "默认分类"
			} else {
				cate = categories_map[cateID].Name
			}
			title_articles[cate] = append(title_articles[cate], article)
		}

		c.HTML(http.StatusOK, "blog/timeline.html", pongo2.Context{
			"site_info":      site_info,
			"page_title":     title,
			"person_info":    person_info,
			"title_articles": title_articles,
		})
	}
}

func (controller *BlogController) Journals(c *gin.Context) {
	if c.Request.Method == "GET" {
		site_info := common.GenerateSiteInfo()
		person_info := common.GeneratePersonInfo()
		title := "日志"
		var (
			journals []*models.Journal
			err      error
		)
		if journals, err = models.NewJournal().FindAll(); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("find journals error: %v", err))
		}

		c.HTML(http.StatusOK, "blog/journals.html", pongo2.Context{
			"site_info":   site_info,
			"page_title":  title,
			"person_info": person_info,
			"journals":    journals,
		})
	}
}
