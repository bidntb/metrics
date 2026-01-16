package middleware

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandler(notFoundHandler, badRequestHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypeBind {
					badRequestHandler(c)
					c.Abort()
					return
				}
			}
			notFoundHandler(c)
			c.Abort()
		}
	}
}
