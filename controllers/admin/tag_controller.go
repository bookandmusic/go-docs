package admin

import (
	"net/http"
	"strconv"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type TagController struct{}

func NewTagController() *TagController {
	return &TagController{}
}

func (tag *TagController) Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/tag/admin.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
		})
	}
}

func (u *TagController) TagList(c *gin.Context) {
	keyword := c.Query("keyword")
	var categories []*models.Tag
	if keyword != "" {
		categories, _ = models.NewTag().FindByKeyword(keyword)
	} else {
		categories, _ = models.NewTag().FindAll()
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": len(categories), "data": categories})
	return
}

func (u *TagController) DeleteTag(c *gin.Context) {
	var tags []*models.Tag
	if err := c.ShouldBindJSON(&tags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "无效的参数"})
		return
	}
	var tagIds []int

	for _, cate := range tags {
		tagIds = append(tagIds, int(cate.ID))
	}
	num, err := models.NewTag().DeleteByTagIds(tagIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除标签失败"})
		return
	}
	if num == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "标签不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除标签成功"})
	return
}

func (tag *TagController) EditTag(c *gin.Context) {
	// 从路由中获取整数参数
	var (
		tagID int
		err   error
		obj   *models.Tag
	)
	tagIDStr := c.Query("id")
	if tagIDStr == "" || tagIDStr == "0" {
		obj = models.NewTag()
	} else {
		// 将字符串形式的 tagID 转换为整数
		tagID, err = strconv.Atoi(tagIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
			return
		}
		obj, _ = models.NewTag().FindByTagId(tagID)
	}

	if c.Request.Method == "GET" {
		// 处理 GET 请求，返回登录页面
		c.HTML(http.StatusOK, "admin/tag/edit.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
			"tag":       obj,
		})
	} else if c.Request.Method == "POST" {
		name := c.PostForm("name")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "名称是必须的"})
			return
		}
		if obj != nil && name == obj.Name {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "名称没有任何改变"})
			return
		}
		exists, _ := models.NewTag().FindByName(name)
		if exists != nil && exists.ID != 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "标签名已存在"})
			return
		}

		if obj == nil || obj.ID == 0 {
			if _, err := models.NewTag().Create(name); err != nil {
				global.GVA_LOG.Warn("Failed to add tag", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "编辑标签失败"})
				return
			}
		} else {
			updates := map[string]interface{}{
				"name": name,
			}
			if err := obj.Update(updates); err != nil {
				global.GVA_LOG.Warn("Failed to edit tag", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "编辑标签失败"})
				return
			}
		}
		return
	}
}
