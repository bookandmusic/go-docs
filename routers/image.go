package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/controllers/image"
)

func InitImageRoutes(router *gin.Engine) {
	imageController := image.NewImageController()

	imageRouter := router.Group("/image")
	{
		imageRouter.GET("/random", imageController.RandomImg)
	}
}
