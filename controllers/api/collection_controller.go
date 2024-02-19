package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type CollectionAPIController struct{}

type CollectionJson struct {
	ID     int    `json:"id"`
	Name   string `json:"name" binding:"required"`
	Author string `json:"author" binding:"required"`
}

func NewCollectionAPIController() *CollectionAPIController {
	return &CollectionAPIController{}
}

func (controller *CollectionAPIController) CollectionList(c *gin.Context) {
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	var collections []*models.Collection

	collections, _ = models.NewCollection().FindByKeyword(keyword, sort, true, false)

	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: collections})
}

func (controller *CollectionAPIController) DeleteCollection(c *gin.Context) {
	collectionIdStr := c.Param("collectionId")
	collectionId, err := strconv.Atoi(collectionIdStr)
	if err != nil || collectionId <= 0 {
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	var collectionIds []int

	collectionIds = append(collectionIds, collectionId)

	if articleCount := models.NewCollection().ArticleCountByCollectionIds(collectionIds); articleCount != 0 {
		c.JSON(http.StatusOK, common.NotEmptyError)
		return
	}

	num, err := models.NewCollection().DeleteByCollectionIds(collectionIds)
	if err != nil {
		c.JSON(http.StatusOK, common.DeleteError)
		return
	}
	if num == 0 {
		c.JSON(http.StatusOK, common.NotExistsError)
		return
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: ""})
}

func (controller *CollectionAPIController) EditCollection(c *gin.Context) {
	// 1. gin中获取json中的username和password
	var (
		json CollectionJson
		obj  *models.Collection
	)

	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	name := json.Name
	author := json.Author
	collectionID := json.ID

	if collectionID != -1 {
		obj, _ = models.NewCollection().FindByCollectionId(collectionID)
		if obj != nil && name == obj.Name {
			c.JSON(http.StatusOK, common.NoUpdateError)
			return
		}
	}

	exists, _ := models.NewCollection().FindByName(name)
	if exists != nil && exists.ID != 0 {
		c.JSON(http.StatusOK, common.ExistsError)
		return
	}

	if obj == nil || obj.ID == 0 {
		obj, err := models.NewCollection().Create(name, author)
		if err != nil {
			global.GVA_LOG.Warn("Failed to add collection", err)
			c.JSON(http.StatusOK, common.CreateError)
			return
		}
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: obj})
		return
	} else {
		updates := map[string]interface{}{
			"name":   name,
			"author": author,
		}
		if err := models.NewCollection().Update(obj, updates); err != nil {
			global.GVA_LOG.Warn("Failed to edit collection", err)
			c.JSON(http.StatusOK, common.UpdateError)
			return
		}
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: obj})
		return
	}
}

func (controller *CollectionAPIController) TocList(c *gin.Context) {
	// 从路由中获取整数参数
	collectionIDStr := c.Param("collectionId")

	// 将字符串形式的 collectionID 转换为整数
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	obj := models.NewArticle()
	articles, _ := obj.FindByCollectionId(collectionID)

	tocList := obj.TocList(articles)
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: tocList})

}

func (controller *CollectionAPIController) TocSorted(c *gin.Context) {

	var sort_data []*models.CollectionTocItem
	if err := c.BindJSON(&sort_data); err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("params error: %v", err))
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: ""})
	if len(sort_data) > 0 {
		if article, err := models.NewArticle().FindByArticleId(sort_data[0].ID); err == nil {
			if err := article.Collection.Update(&article.Collection, map[string]interface{}{"first_doc": article.Identify}); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("update collection %s first doc error: %v", article.Collection.Name, err))
			}
		}
	}
	models.NewArticle().UpdateArticleSort(sort_data)

}

func (controller *CollectionAPIController) BatchDeleteCollection(c *gin.Context) {
	type collectionIds struct {
		IDs []int
	}

	var json collectionIds
	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.ParamError)
		return
	}

	if articleCount := models.NewCollection().ArticleCountByCollectionIds(json.IDs); articleCount != 0 {
		c.JSON(http.StatusOK, common.NotEmptyError)
		return
	}

	num, err := models.NewCollection().DeleteByCollectionIds(json.IDs)
	if err != nil {
		c.JSON(http.StatusOK, common.DeleteError)
		return
	}
	if num == 0 {
		c.JSON(http.StatusOK, common.NotExistsError)
		return
	}
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: ""})
}
