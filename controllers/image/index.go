package image

import (
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/common"
	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

type ImageController struct{}

func NewImageController() *ImageController {
	return &ImageController{}
}

func (controller ImageController) RandomImg(c *gin.Context) {
	category := c.Query("category")
	var imageDir string
	switch category {
	case "cover":
		imageDir = path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "static/image/cover")
	default:
		imageDir = path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "static/image/cover")
	}

	var imagePaths []string

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 判断是否为文件，并且是图片文件
		if !info.IsDir() && utils.IsImageFile(path) {
			imagePaths = append(imagePaths, path)
		}
		return nil
	})

	if err != nil || len(imagePaths) == 0 {
		c.JSON(http.StatusOK, common.SuccessMsg{Code: common.SuccessCode, Data: map[string]any{"image_paths": imagePaths}})
		return
	}
	// 从图片路径中随机选择一个
	imageIndex := rand.Intn(len(imagePaths))
	imagePath := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, imagePaths[imageIndex])

	fileBytes, _ := os.ReadFile(imagePath)

	// 直接将文件内容写入响应体
	c.Writer.WriteString(string(fileBytes))
}
