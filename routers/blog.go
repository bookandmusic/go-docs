package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/controllers/blog"
)

func InitBlogRoutes(router *gin.Engine) {
	blogController := blog.NewBlogController()

	blogRouter := router.Group("/")
	{
		blogRouter.GET("/", blogController.Index)
		blogRouter.GET("/page/:page", blogController.Index)
		blogRouter.GET("/tags/:tag", blogController.Index)
		blogRouter.GET("/categories/:category", blogController.Index)
		blogRouter.GET("/archives", blogController.Archives)
		blogRouter.GET("/archives/:identify", blogController.Post)
		blogRouter.GET("/categories", blogController.Categories)
		blogRouter.GET("/collections", blogController.Collections)
		blogRouter.GET("/journals", blogController.Journals)
	}
}
