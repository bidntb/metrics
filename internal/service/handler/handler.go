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
	mtype := c.Param("type")
	name := c.Param("name")
	valueStr := c.Param("value")

	if mtype != "" && name != "" && valueStr != "" {
		if mtype == "gauge" {
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value for gauge"})
				return
			}
			resp, err := h.svc.UpdateGauge(name, value)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		} else if mtype == "counter" {
			delta, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value for counter"})
				return
			}
			resp, err := h.svc.UpdateCounter(name, delta)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid metric type"})
		return
	}

	var req metrics.UpdateMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	if req.ID == "" || req.MType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id or type"})
		return
	}

	switch req.MType {
	case "gauge":
		if req.Value == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "value required for gauge"})
			return
		}
		resp, err := h.svc.UpdateGauge(req.ID, *req.Value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return

	case "counter":
		if req.Delta == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "delta required for counter"})
			return
		}
		resp, err := h.svc.UpdateCounter(req.ID, *req.Delta)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type: must be gauge or counter"})
		return
	}
}

func (h *Handler) GetValue(c *gin.Context) {
	mtype := c.Param("type")
	name := c.Param("name")

	if mtype == "" || name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing type or name"})
		return
	}

	val, ok := h.svc.GetMetricValue(mtype, name)
	if val == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.String(http.StatusOK, val)
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
