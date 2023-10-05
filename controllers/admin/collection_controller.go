package admin

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
	"github.com/bookandmusic/docs/utils"
)

type CollectionController struct{}

func NewCollectionController() *CollectionController {
	return &CollectionController{}
}

func (ccontroller *CollectionController) Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/collection/admin.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
		})
	}
}

func (controller *CollectionController) CollectionList(c *gin.Context) {
	keyword := c.Query("keyword")
	var categories []*models.Collection

	categories, _ = models.NewCollection().FindByKeyword(keyword, true, false)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": len(categories), "data": categories})
	return
}

func (controller *CollectionController) DeleteCollection(c *gin.Context) {
	var collections []*models.Collection
	if err := c.ShouldBindJSON(&collections); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "无效的参数"})
		return
	}
	var collectionIds []int

	for _, cate := range collections {
		collectionIds = append(collectionIds, int(cate.ID))
	}
	num, err := models.NewCollection().DeleteByCollectionIds(collectionIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除文集失败"})
		return
	}
	if num == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "文集不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除文集成功"})
	return
}

func (controller *CollectionController) EditCollection(c *gin.Context) {
	// 从路由中获取整数参数
	var obj *models.Collection

	collectionIDStr := c.Query("id")
	if collectionIDStr == "" || collectionIDStr == "0" {
		obj = models.NewCollection()
	} else {
		// 将字符串形式的 collectionID 转换为整数
		collectionID, err := strconv.Atoi(collectionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文集ID"})
			return
		}
		obj, _ = models.NewCollection().FindByCollectionId(collectionID)
	}

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/collection/edit.html", pongo2.Context{
			"csrfToken":  csrf.Token(c.Request),
			"collection": obj,
		})
	} else if c.Request.Method == "POST" {
		name := c.PostForm("name")
		author := c.PostForm("author")
		cover := c.PostForm("cover")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "名称是必须的"})
			return
		}

		if obj != nil && name != obj.Name {
			exists, _ := models.NewCollection().FindByName(name)
			if exists != nil && exists.ID != 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "文集已存在"})
				return
			}
		}

		if cover == "" {
			originCoverPath := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "static/public/images/cover.jpeg")
			outCoverPath := path.Join(global.GVA_CONFIG.Server.UploadPath, utils.GenerateMD5Hash(name)+".png")
			fontPath := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "static/public/fonts/songti.ttf")
			if err := utils.GenerateCover(originCoverPath, fontPath, name, outCoverPath); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("generate cover error: %v", err))
				cover = originCoverPath
			} else {
				cover = outCoverPath
			}
			// 判断路径是否以 "/" 开头
			if !strings.HasPrefix(cover, "/") {
				// 如果不是以 "/" 开头，添加 "/"
				cover = "/" + cover
			}
		}

		if obj == nil || obj.ID == 0 {
			if err := models.NewCollection().Create(name, author, cover); err != nil {
				global.GVA_LOG.Warn("Failed to add collection", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "编辑文集失败"})
				return
			}
		} else {
			updates := map[string]interface{}{
				"name":   name,
				"author": author,
				"cover":  cover,
			}
			if err := obj.Update(updates); err != nil {
				global.GVA_LOG.Warn("Failed to edit collection", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "编辑文集失败"})
				return
			}
		}

		return
	}
}

func (controller *CollectionController) TocList(c *gin.Context) {
	// 从路由中获取整数参数
	collectionIDStr := c.Query("collection_id")

	// 将字符串形式的 collectionID 转换为整数
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文集ID"})
		return
	}
	collection, _ := models.NewCollection().FindByCollectionId(collectionID)
	if c.Request.Method == "GET" {
		obj := models.NewArticle()
		articles, _ := obj.FindByCollectionId(collectionID)

		tocList := obj.TocList(articles)
		c.HTML(http.StatusOK, "admin/collection/toc_list.html", pongo2.Context{
			"csrfToken":  csrf.Token(c.Request),
			"toc_list":   tocList,
			"collection": collection,
		})
	} else if c.Request.Method == "POST" {
		var sort_data []*models.CollectionTocItem
		if err := c.BindJSON(&sort_data); err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("params error: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "参数有误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "msg": ""})
		if len(sort_data) > 0 {
			if article, err := models.NewArticle().FindByArticleId(sort_data[0].ID); err == nil {
				if err := article.Collection.Update(map[string]interface{}{"first_doc": article.Identify}); err != nil {
					global.GVA_LOG.Error(fmt.Sprintf("update collection %s first doc error: %v", article.Collection.Name, err))
				}
			}
		}
		models.NewArticle().UpdateArticleSort(sort_data)
		return
	}

	return
}
