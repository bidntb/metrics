package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/handler"
	"bidntb/metrics/internal/nconfig"
)

func main() {

	ServerAddress := nconfig.GetServerAddress()

	router := gin.Default()

	router.GET("/", handler.IndexHandler)
	router.GET("/value/:type/:name", handler.ValueHandler)

	router.POST("/update/counter/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/counter", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/gauge/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/gauge", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/:wrong", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	})
	router.POST("/update/:wrong/*any", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	})

	router.POST("/update/gauge/:name/:value", handler.AddGaugeHandler)
	router.POST("/update/counter/:name/:value", handler.AddCounterHandler)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	if err := router.Run(ServerAddress); err != nil {
		panic(err)
	}
}
