package app

import (
	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/handler"
	"bidntb/metrics/internal/logger"
	"bidntb/metrics/internal/nconfig"
	"bidntb/metrics/internal/storage"
)

func Run() {

	serverAddress := nconfig.GetServerAddress()

	router := gin.Default()
	router.Use(logger.LoggingMiddleware())

	storageInstance := storage.NewMemStorage()

	h := handler.NewHandler(storageInstance)
	router.GET("/", h.IndexHandler)
	router.GET("/value/:type/:name", h.ValueHandler)

	router.POST("/update/counter/", h.NotFoundHandler)
	router.POST("/update/counter", h.NotFoundHandler)
	router.POST("/update/gauge/", h.NotFoundHandler)
	router.POST("/update/gauge", h.NotFoundHandler)
	router.POST("/update/:wrong", h.BadRequestHandler)
	router.POST("/update/:wrong/*any", h.BadRequestHandler)

	router.POST("/update/gauge/:name/:value", h.AddGaugeHandler)
	router.POST("/update/counter/:name/:value", h.AddCounterHandler)

	router.NoRoute(h.NotFoundHandler)

	if err := router.Run(serverAddress); err != nil {
		panic(err)
	}
}
