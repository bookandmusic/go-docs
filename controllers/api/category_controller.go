package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type CategoryAPIController struct{}

type CategoryJson struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

func NewCategoryAPIController() *CategoryAPIController {
	return &CategoryAPIController{}
}

func (controller *CategoryAPIController) CategoryList(c *gin.Context) {
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	var categories []*models.Category
	categories, _ = models.NewCategory().FindByKeyword(keyword, sort)

	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: categories})
}

func (controller *CategoryAPIController) DeleteCategory(c *gin.Context) {
	categoryIdStr := c.Param("categoryId")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId <= 0 {
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	var categoryIds []int

	categoryIds = append(categoryIds, categoryId)
	if articleCount := models.NewCategory().ArticleCountByCategoryIds(categoryIds); articleCount != 0 {
		c.JSON(http.StatusOK, common.NotEmptyError)
		return
	}

	num, err := models.NewCategory().DeleteByCategoryIds(categoryIds)
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

func (controller *CategoryAPIController) EditCategory(c *gin.Context) {
	// 1. gin中获取json中的username和password
	var (
		json CategoryJson
		obj  *models.Category
	)

	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	name := json.Name
	categoryID := json.ID

	if categoryID != -1 {
		obj, _ = models.NewCategory().FindByCategoryId(categoryID)
		if obj != nil && name == obj.Name {
			c.JSON(http.StatusOK, common.NoUpdateError)
			return
		}
	}

	exists, _ := models.NewCategory().FindByName(name)
	if exists != nil && exists.ID != 0 {
		c.JSON(http.StatusOK, common.ExistsError)
		return
	}

	if obj == nil || obj.ID == 0 {
		obj, err := models.NewCategory().Create(name)
		if err != nil {
			global.GVA_LOG.Warn("Failed to add category", err)
			c.JSON(http.StatusOK, common.CreateError)
			return
		}
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: obj})
		return
	} else {
		updates := map[string]interface{}{
			"name": name,
		}
		if err := models.NewCategory().Update(obj, updates); err != nil {
			global.GVA_LOG.Warn("Failed to edit category", err)
			c.JSON(http.StatusOK, common.UpdateError)
			return
		}
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: obj})
		return
	}
}

func (controller *CategoryAPIController) BatchDeleteCategory(c *gin.Context) {
	type categoryIds struct {
		IDs []int
	}

	var json categoryIds
	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.ParamError)
		return
	}

	if articleCount := models.NewCategory().ArticleCountByCategoryIds(json.IDs); articleCount != 0 {
		c.JSON(http.StatusOK, common.NotEmptyError)
		return
	}

	num, err := models.NewCategory().DeleteByCategoryIds(json.IDs)
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
