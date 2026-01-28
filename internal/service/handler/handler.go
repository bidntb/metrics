package handler

import (
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
	if mtype, name, valueStr := c.Param("type"), c.Param("name"), c.Param("value"); mtype != "" && name != "" && valueStr != "" {
		if resp, err := h.updateFromPath(c, mtype, name, valueStr); err != nil {
			return
		} else {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	var req metrics.UpdateMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	resp, err := h.svc.UpdateMetric(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetValue(c *gin.Context) {
	mtype, name := c.Param("type"), c.Param("name")
	if mtype == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing type or name"})
		return
	}

	if val, ok := h.svc.GetMetricValue(mtype, name); ok && val != "" {
		c.String(http.StatusOK, val)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func (h *Handler) GetMetricJSON(c *gin.Context) {
	var req metrics.GetMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	resp, err := h.svc.GetMetric(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if resp.Value == "" {
		c.JSON(http.StatusOK, gin.H{
			"id":    resp.ID,
			"MType": resp.MType,
			"Value": nil,
			"Delta": nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListMetrics(c *gin.Context) {
	metricsMap := h.svc.ListAll()
	c.JSON(http.StatusOK, gin.H{
		"gauge":   metricsMap["gauge"],
		"counter": metricsMap["counter"],
	})
}

func (h *Handler) NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
}

func (h *Handler) BadRequestHandler(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
}

func (h *Handler) updateFromPath(c *gin.Context, mtype, name, valueStr string) (*metrics.MetricResponse, error) {
	switch mtype {
	case "gauge":
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value for gauge"})
			return nil, err
		}
		return h.svc.UpdateGauge(name, value)
	case "counter":
		delta, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value for counter"})
			return nil, err
		}
		return h.svc.UpdateCounter(name, delta)
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid metric type"})
		return nil, nil
	}
}
