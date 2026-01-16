package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/service/metrics"
)

type Handler struct {
	svc *metrics.Service
}

func NewHandler(svc *metrics.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) UpdateMetric(c *gin.Context) {
	var req metrics.UpdateMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.UpdateMetric(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateGauge(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing name"})
		return
	}

	var req metrics.UpdateGaugeRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		h.svc.UpdateGauge(name, req.Value)
		c.Status(http.StatusOK)
		return
	}

	valueStr := c.Param("value")
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		c.Error(err)
		return
	}
	h.svc.UpdateGauge(name, value)
	c.Status(http.StatusOK)
}

func (h *Handler) UpdateCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing name"})
		return
	}

	var req metrics.UpdateCounterRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		h.svc.UpdateCounter(name, req.Delta)
		c.Status(http.StatusOK)
		return
	}

	valueStr := c.Param("value")
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	h.svc.UpdateCounter(name, value)
	c.Status(http.StatusOK)
}

func (h *Handler) GetValue(c *gin.Context) {
	mtype := c.Param("type")
	name := c.Param("name")
	if mtype == "" || name == "" {
		c.Error(fmt.Errorf("missing params"))
		return
	}
	if val, ok := h.svc.GetMetricValue(mtype, name); ok {
		c.String(http.StatusOK, val)
		return
	}
	c.Status(http.StatusNotFound)
}

func (h *Handler) ListMetrics(c *gin.Context) {
	metrics := h.svc.ListAll()
	c.JSON(http.StatusOK, gin.H{
		"gauge":   metrics["gauge"],
		"counter": metrics["counter"],
	})
}

func (h *Handler) NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
}

func (h *Handler) BadRequestHandler(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
}

func (h *Handler) GetMetricJSON(c *gin.Context) {
	var req metrics.GetMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if resp, found := h.svc.GetMetric(req); found {
		c.JSON(http.StatusOK, resp)
		return
	}
	c.Status(http.StatusNotFound)
}
