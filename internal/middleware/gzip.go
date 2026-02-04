package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(strings.ToLower(c.GetHeader("Content-Encoding")), "gzip") {
			r, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid gzip"})
				return
			}
			body, _ := io.ReadAll(r)
			err = r.Close()
			if err != nil {
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Request.Header.Del("Content-Encoding")
		}

		ae := c.GetHeader("Accept-Encoding")
		if strings.Contains(strings.ToLower(ae), "gzip") {
			c.Writer.Header().Set("Content-Encoding", "gzip")
			c.Writer.Header().Add("Vary", "Accept-Encoding")
			w := gzip.NewWriter(c.Writer)
			c.Writer = &gzipWriter{ResponseWriter: c.Writer, w: w}
			defer func(w *gzip.Writer) {
				err := w.Close()
				if err != nil {

				}
			}(w)
		}

		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	w *gzip.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.w.Write(data)
}
