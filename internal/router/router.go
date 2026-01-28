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

	r.POST("/update/:type/:name/:value", h.UpdateMetric)
	r.POST("/update/", h.UpdateMetric)
	r.GET("/value/:type/:name", h.GetValue)
	r.POST("/value/", h.GetMetricJSON)
	r.GET("/", h.ListMetrics)

	r.POST("/update/counter", h.NotFoundHandler)
	r.POST("/update/gauge/", h.NotFoundHandler)
	r.POST("/update/gauge", h.NotFoundHandler)
	r.NoRoute(h.NotFoundHandler)

	return r
}
