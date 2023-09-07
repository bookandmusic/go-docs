package middlewares

import (
	"github.com/gin-gonic/gin"
	minify "github.com/tdewolff/minify/v2"
)

type MinifyResponseWriter struct {
	gin.ResponseWriter
	m *minify.M
}

func (mrw *MinifyResponseWriter) Write(data []byte) (int, error) {
	// 在这里进行响应数据的最小化
	ct := mrw.Header().Get("Content-Type")
	minData, err := mrw.m.Bytes(ct, data)
	if err != nil {
		return 0, err
	}
	return mrw.ResponseWriter.Write(minData)
}

func MinifyMiddleware(m *minify.M) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用自定义 MinifyResponseWriter 替换原始的 ResponseWriter
		writer := c.Writer
		minifyWriter := &MinifyResponseWriter{ResponseWriter: writer, m: m}
		c.Writer = minifyWriter

		// 继续处理请求
		c.Next()
	}
}
