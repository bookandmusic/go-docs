package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
)

type DashboardAPIController struct{}

func NewDashboardAPIController() *DashboardAPIController {
	return &DashboardAPIController{}
}

func (controller *DashboardAPIController) Index(c *gin.Context) {
	article_info := common.GenerateArticleInfo()
	c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: article_info})
}
