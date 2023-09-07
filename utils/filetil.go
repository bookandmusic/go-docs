package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func PathExists(path string) error {
	// 检查路径是否存在，如果不存在则创建路径
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
	}
	return nil
}

func GenerateFileNameByMD5(file *multipart.FileHeader) string {
	// 打开上传的文件
	openedFile, err := file.Open()
	if err != nil {
		return ""
	}
	defer openedFile.Close()

	// 计算文件内容的MD5哈希值
	hash := md5.New()
	if _, err := io.Copy(hash, openedFile); err != nil {
		return ""
	}

	// 将MD5哈希值转换为16进制字符串作为文件名
	md5Hash := hex.EncodeToString(hash.Sum(nil))
	return md5Hash
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

// 判断指定目录下是否存在指定后缀的文件
func HasFileOfExt(path string, exts []string) bool {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			ext := filepath.Ext(info.Name())

			for _, item := range exts {
				if strings.EqualFold(ext, item) {
					return os.ErrExist
				}
			}

		}
		return nil
	})

	return err == os.ErrExist
}

// 拷贝文件
func CopyFile(source string, dst string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	_, err = os.Stat(filepath.Dir(dst))

	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(dst), 0o766)
		} else {
			return err
		}
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dst, sourceInfo.Mode())
		}

	}

	return
}

// 拷贝目录
func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourceFilePointer := filepath.Join(source, obj.Name())

		destinationFilePointer := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourceFilePointer, destinationFilePointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourceFilePointer, destinationFilePointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

// 忽略字符串中的BOM头
func ReadFileAndIgnoreUTF8BOM(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	data = bytes.Replace(data, []byte("\r"), []byte(""), -1)
	if len(data) >= 3 && data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
		return data[3:], err
	}

	return data, nil
}

func EscapeJS(md string) string {
	// Replace backslashes
	md = strings.ReplaceAll(md, "\\", "\\\\")

	// Replace single quotes
	md = strings.ReplaceAll(md, "'", "\\'")

	// Replace double quotes
	md = strings.ReplaceAll(md, "\"", "\\\"")

	// Replace newlines
	md = strings.ReplaceAll(md, "\n", "\\n")

	// Replace carriage returns
	md = strings.ReplaceAll(md, "\r", "\\r")

	// Replace open and close parentheses
	md = strings.ReplaceAll(md, "(", "\\(")
	md = strings.ReplaceAll(md, ")", "\\)")

	// Replace plus sign
	md = strings.ReplaceAll(md, "+", "\\+")

	// Replace question mark
	md = strings.ReplaceAll(md, "?", "\\?")

	// Replace square brackets
	md = strings.ReplaceAll(md, "[", "\\[")
	md = strings.ReplaceAll(md, "]", "\\]")

	// Replace dollar sign
	md = strings.ReplaceAll(md, "$", "\\$")

	// Replace caret
	md = strings.ReplaceAll(md, "^", "\\^")
	return md
}

func MarkdownToHTML(markdownStr string) (string, error) {
	// 自定义解析器
	var (
		htmlBuffer  bytes.Buffer
		htmlContent string
		err         error
	)
	markdown := goldmark.New(
		// 支持 GFM
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	err = markdown.Convert([]byte(markdownStr), &htmlBuffer)
	if err == nil {
		htmlContent = htmlBuffer.String()
	}

	return htmlContent, err
}

func GenerateCover(originImgPath, titleFontPath, title, outImgPath string) error {
	// 打开原始图片
	file, err := os.Open(originImgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取图片
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 调整图片大小（如果需要的话）
	img = resize.Resize(600, 800, img, resize.Lanczos3)

	// 创建一个画布
	dc := gg.NewContext(600, 800)

	// 将图片绘制到画布
	dc.DrawImage(img, 0, 0)

	// 设置字体和标题
	err = dc.LoadFontFace(titleFontPath, 36)
	if err != nil {
		return err
	}

	// 获取文本的宽度
	width, _ := dc.MeasureString(title)

	// 计算居中的 x 坐标
	x := width/2 + 300
	y := 200

	// 设置字体颜色
	dc.SetColor(color.Black)

	// 添加标题
	dc.DrawStringAnchored(title, float64(x), float64(y), 1, 1)

	// 保存图片
	if err := dc.SavePNG(outImgPath); err != nil {
		return err
	}

	return nil
}
