package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
)

type TagAPIController struct{}

type TagJson struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

func NewTagAPIController() *TagAPIController {
	return &TagAPIController{}
}

func (controller *TagAPIController) TagList(c *gin.Context) {
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	var tags []*models.Tag

	tags, _ = models.NewTag().FindByKeyword(keyword, sort)

	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: tags})
}

func (controller *TagAPIController) DeleteTag(c *gin.Context) {
	tagIdStr := c.Param("tagId")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil || tagId <= 0 {
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	var tagIds []int

	tagIds = append(tagIds, tagId)

	if articleCount := models.NewTag().ArticleCountByTagIds(tagIds); articleCount != 0 {
		c.JSON(http.StatusOK, common.NotEmptyError)
		return
	}

	num, err := models.NewTag().DeleteByTagIds(tagIds)
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

func (controller *TagAPIController) EditTag(c *gin.Context) {
	// 1. gin中获取json中的username和password
	var (
		json TagJson
		obj  *models.Tag
	)

	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.ParamError)
		return
	}
	name := json.Name
	tagID := json.ID

	if tagID != -1 {
		obj, _ = models.NewTag().FindByTagId(tagID)
		if obj != nil && name == obj.Name {
			c.JSON(http.StatusOK, common.NoUpdateError)
			return
		}
	}

	exists, _ := models.NewTag().FindByName(name)
	if exists != nil && exists.ID != 0 {
		c.JSON(http.StatusOK, common.ExistsError)
		return
	}

	if obj == nil || obj.ID == 0 {
		obj, err := models.NewTag().Create(name)
		if err != nil {
			global.GVA_LOG.Warn("Failed to add tag: ", err)
			c.JSON(http.StatusOK, common.CreateError)
			return
		}
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: obj})
		return
	} else {
		updates := map[string]interface{}{
			"name": name,
		}
		if err := models.NewTag().Update(obj, updates); err != nil {
			global.GVA_LOG.Warn("Failed to edit tag: ", err)
			c.JSON(http.StatusOK, common.UpdateError)
			return
		}
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: obj})
		return
	}
}

func (controller *TagAPIController) BatchDeleteTag(c *gin.Context) {
	type tagIds struct {
		IDs []int
	}

	var json tagIds
	if err := c.ShouldBindJSON(&json); err != nil {
		// 如果绑定失败，返回错误信息
		c.JSON(http.StatusOK, common.ParamError)
		return
	}

	if articleCount := models.NewTag().ArticleCountByTagIds(json.IDs); articleCount != 0 {
		c.JSON(http.StatusOK, common.NotEmptyError)
		return
	}

	num, err := models.NewTag().DeleteByTagIds(json.IDs)
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
