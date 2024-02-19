package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/controllers/api"
)

func InitAPIRoutes(router *gin.Engine) {
	userAPIController := api.NewUserAPIController()
	tagAPIController := api.NewTagAPIController()
	categoryAPIController := api.NewCategoryAPIController()
	collectionAPIController := api.NewCollectionAPIController()
	articleAPIController := api.NewArticleAPIController()
	dashboardAPIController := api.NewDashboardAPIController()

	apiRouter := router.Group(common.APIPrefix)
	{
		apiRouter.GET(common.APIDashboard, dashboardAPIController.Index)
		apiRouter.POST(common.APILoginUrl, userAPIController.Login)
		apiRouter.POST(common.APILogoutUrl, userAPIController.Logout)
		apiRouter.GET(common.APIUserInfoUrl, userAPIController.UserInfo)

		apiRouter.GET(common.APITagUrl, tagAPIController.TagList)
		apiRouter.POST(common.APITagUrl, tagAPIController.EditTag)
		apiRouter.DELETE(common.APITagUrl, tagAPIController.BatchDeleteTag)
		apiRouter.DELETE(common.APITagDetailUrl, tagAPIController.DeleteTag)

		apiRouter.GET(common.APICategoryUrl, categoryAPIController.CategoryList)
		apiRouter.POST(common.APICategoryUrl, categoryAPIController.EditCategory)
		apiRouter.DELETE(common.APICategoryUrl, categoryAPIController.BatchDeleteCategory)
		apiRouter.DELETE(common.APICategoryDetailUrl, categoryAPIController.DeleteCategory)

		apiRouter.GET(common.APICollectionUrl, collectionAPIController.CollectionList)
		apiRouter.POST(common.APICollectionUrl, collectionAPIController.EditCollection)
		apiRouter.DELETE(common.APICollectionUrl, collectionAPIController.BatchDeleteCollection)
		apiRouter.DELETE(common.APICollectionDetailUrl, collectionAPIController.DeleteCollection)

		apiRouter.GET(common.APIArticleUrl, articleAPIController.ArticleList)
		apiRouter.GET(common.APIArticleDetailUrl, articleAPIController.ArticleDetail)
		apiRouter.POST(common.APIArticleUrl, articleAPIController.EditArticle)
		apiRouter.DELETE(common.APIArticleUrl, articleAPIController.BatchDeleteArticle)
		apiRouter.DELETE(common.APIArticleDetailUrl, articleAPIController.DeleteArticle)
		apiRouter.DELETE(common.APIArticImportleUrl, articleAPIController.ImportMds)
	}
}
