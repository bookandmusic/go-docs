package admin

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (u *UploadController) UploadImg(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "上传失败",
			"url":    "",
		})
		return
	}

	// 生成文件名
	filename := utils.GenerateFileNameByMD5(file)
	imgName := filename + ".png"
	imgPath := filepath.Join(global.GVA_CONFIG.Server.UploadPath, imgName)
	// 保存文件
	if err := c.SaveUploadedFile(file, imgPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"msg":    "文件保存失败",
			"url":    "",
		})
		return
	}

	// 判断路径是否以 "/" 开头
	if !strings.HasPrefix(imgPath, "/") {
		// 如果不是以 "/" 开头，添加 "/"
		imgPath = "/" + imgPath
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "上传成功",
		"url":    imgPath,
	})
}
