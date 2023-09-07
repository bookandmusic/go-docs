package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/controllers/admin"
)

func InitAdminRoutes(router *gin.Engine) {
	userController := admin.NewUserController()
	adminController := admin.NewAdminController()
	categoryController := admin.NewCategoryController()
	tagController := admin.NewTagController()
	collectionController := admin.NewCollectionController()
	articleController := admin.NewArticleController()
	uploadController := admin.NewUploadController()
	journalController := admin.NewJournalController()

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/", adminController.Index)
		adminRouter.GET("/welcome", adminController.Welcome)

		adminRouter.POST("/login", userController.Login)
		adminRouter.GET("/login", userController.Login)
		adminRouter.GET("/logout", userController.Logout)
		adminRouter.GET("/user-setting", userController.UserSetting)
		adminRouter.POST("/user-setting", userController.UserSetting)
		adminRouter.GET("/user-password", userController.UserPassword)
		adminRouter.POST("/user-password", userController.UserPassword)

		adminRouter.GET("/category/edit", categoryController.EditCategory)
		adminRouter.POST("/category/edit", categoryController.EditCategory)
		adminRouter.POST("/category/delete", categoryController.DeleteCategory)
		adminRouter.GET("/category", categoryController.Index)

		adminRouter.GET("/tag/edit", tagController.EditTag)
		adminRouter.POST("/tag/edit", tagController.EditTag)
		adminRouter.POST("/tag/delete", tagController.DeleteTag)
		adminRouter.GET("/tag", tagController.Index)

		adminRouter.GET("/collection/edit", collectionController.EditCollection)
		adminRouter.POST("/collection/edit", collectionController.EditCollection)
		adminRouter.POST("/collection/delete", collectionController.DeleteCollection)
		adminRouter.GET("/collection/toclist", collectionController.TocList)
		adminRouter.POST("/collection/toclist", collectionController.TocList)
		adminRouter.GET("/collection", collectionController.Index)

		adminRouter.GET("/article/edit", articleController.EditArticle)
		adminRouter.POST("/article/edit", articleController.EditArticle)
		adminRouter.POST("/article/delete", articleController.DeleteArticle)
		adminRouter.GET("/article", articleController.Index)

		adminRouter.GET("/journal", journalController.Index)
		adminRouter.GET("/journal/edit", journalController.EditJournal)
		adminRouter.POST("/journal/edit", journalController.EditJournal)

		adminSystem := adminRouter.Group("/system")
		{
			adminSystem.GET("/import_mds", articleController.ImportMds)
			adminSystem.POST("/import_mds", articleController.ImportMds)
		}

		adminApi := adminRouter.Group("/api")
		{
			adminApi.GET("/menu_data", adminController.AdminCenterMenu)
			adminApi.GET("/category", categoryController.CategoryList)
			adminApi.GET("/tag", tagController.TagList)
			adminApi.GET("/collection", collectionController.CollectionList)
			adminApi.GET("/article", articleController.ArticleList)
			adminApi.GET("/journal", journalController.JournalList)
		}

		uploadApi := adminRouter.Group("/upload")
		{
			uploadApi.POST("/img", uploadController.UploadImg)
		}
	}
}
