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

type CategoryController struct{}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

func (category *CategoryController) Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/category/admin.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
		})
	}
}

func (u *CategoryController) CategoryList(c *gin.Context) {
	keyword := c.Query("keyword")
	var categories []*models.Category
	if keyword != "" {
		categories, _ = models.NewCategory().FindByKeyword(keyword)
	} else {
		categories, _ = models.NewCategory().FindAll()
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": len(categories), "data": categories})
	return
}

func (u *CategoryController) DeleteCategory(c *gin.Context) {
	var categorys []*models.Category
	if err := c.ShouldBindJSON(&categorys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "无效的参数"})
		return
	}
	var categoryIds []int

	for _, cate := range categorys {
		categoryIds = append(categoryIds, int(cate.ID))
	}
	num, err := models.NewCategory().DeleteByCategoryIds(categoryIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除分类失败"})
		return
	}
	if num == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除分类成功"})
	return
}

func (category *CategoryController) EditCategory(c *gin.Context) {
	// 从路由中获取整数参数
	var (
		obj *models.Category
	)
	categoryIDStr := c.Query("id")
	if categoryIDStr == "" || categoryIDStr == "0" {
		obj = models.NewCategory()
	} else {
		// 将字符串形式的 categoryID 转换为整数
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分类ID"})
			return
		}
		obj, _ = models.NewCategory().FindByCategoryId(categoryID)
	}

	if c.Request.Method == "GET" {
		// 处理 GET 请求，返回登录页面
		c.HTML(http.StatusOK, "admin/category/edit.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
			"category":  obj,
		})
	} else if c.Request.Method == "POST" {
		name := c.PostForm("name")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "名称是必须的"})
			return
		}
		if obj != nil && obj.Name == name {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "名称没有任何改变"})
			return
		}
		exists, _ := models.NewCategory().FindByName(name)
		if exists != nil && exists.ID != 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "分类名已存在"})
			return
		}
		if obj == nil || obj.ID == 0 {
			if err := models.NewCategory().Create(name); err != nil {
				global.GVA_LOG.Warn("Failed to add category", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "添加分类失败"})
				return
			}
		} else {
			updates := map[string]interface{}{
				"name": name,
			}
			if err := obj.Update(updates); err != nil {
				global.GVA_LOG.Warn("Failed to edit category", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "更新分类失败"})
				return
			}
		}

		return
	}
}
