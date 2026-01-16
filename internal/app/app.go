package app

import (
	"bidntb/metrics/internal/middleware"
	"bidntb/metrics/internal/service/handler"
	"bidntb/metrics/internal/service/logger"
	"bidntb/metrics/internal/service/metrics"
	"bidntb/metrics/internal/storage"

	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/nconfig"
)

func setupRouter(storage storage.Interface) *gin.Engine {
	metricsSvc := metrics.NewService(storage)
	h := handler.NewHandler(metricsSvc)

	r := gin.New()
	r.Use(logger.LoggingMiddleware())
	r.Use(middleware.ErrorHandler(h.NotFoundHandler, h.BadRequestHandler))

	r.GET("/", h.ListMetrics)
	r.POST("/value", h.GetMetricJSON)
	r.POST("/update", h.UpdateMetric)
	r.POST("/update/gauge/:name/:value", middleware.ValidateGaugeValue(), h.UpdateGauge)
	r.POST("/update/counter/:name/:value", middleware.ValidateCounterValue(), h.UpdateCounter)
	r.GET("/value/:type/:name", h.GetValue)

	r.POST("/update/counter", h.NotFoundHandler)
	r.POST("/update/gauge/", h.NotFoundHandler)
	r.POST("/update/gauge", h.NotFoundHandler)
	r.POST("/update/:wrong", h.BadRequestHandler)
	r.POST("/update/:wrong/*any", h.BadRequestHandler)

	r.NoRoute(h.NotFoundHandler)

	return r
}
func Run() {
	serverAddress := nconfig.GetServerAddress()
	storageInstance := storage.NewMemStorage()

	router := setupRouter(storageInstance)

	if err := router.Run(serverAddress); err != nil {
		panic(err)
	}
}
