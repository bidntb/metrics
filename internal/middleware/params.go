package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ValidateCounterValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Param("value")
		if value == "" || value == "none" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Counter value must be integer"})
			c.Abort()
			return
		}
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Counter value must be integer"})
			c.Abort()
		}
		c.Next()
	}
}

func ValidateGaugeValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Param("value")
		if value == "" || value == "none" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gauge value must be number"})
			c.Abort()
			return
		}

		if _, err := strconv.ParseFloat(value, 64); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gauge value must be number"})
			c.Abort()
		}
		c.Next()
	}
}
