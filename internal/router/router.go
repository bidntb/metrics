package router

import (
	"bidntb/metrics/internal/middleware"
	"bidntb/metrics/internal/service/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandler(h.NotFoundHandler, h.BadRequestHandler))

	r.GET("/", h.ListMetrics)
	r.POST("/value", h.GetMetricJSON)
	r.POST("/value/", h.GetMetricJSON)
	r.POST("/update", h.UpdateMetric)
	r.POST("/update/", h.UpdateMetric)
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
