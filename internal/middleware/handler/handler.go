package handler

import (
	"encoding/json"
	"io"
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
		resp, err := h.updateFromPath(c, mtype, name, valueStr)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var req metrics.UpdateMetricRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	resp, err := h.svc.UpdateMetric(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metricResponsePayload(resp))
}

func metricResponsePayload(resp metrics.MetricResponse) gin.H {
	payload := gin.H{"id": resp.ID, "type": resp.MType}
	if resp.MType == "gauge" {
		if resp.Value != "" {
			v, _ := strconv.ParseFloat(resp.Value, 64)
			payload["value"] = v
		} else {
			payload["value"] = nil
		}
		payload["delta"] = nil
	} else {
		payload["delta"] = resp.Delta
	}
	return payload
}

func (h *Handler) GetValue(c *gin.Context) {
	mtype, name := c.Param("type"), c.Param("name")
	if mtype == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing type or name"})
		return
	}

	if val, ok := h.svc.GetMetricValue(mtype, name); ok && val != "" {
		c.JSON(http.StatusOK, val)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func (h *Handler) GetMetricJSON(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}
	var req metrics.GetMetricRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}
	if req.MType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing type or MType in JSON"})
		return
	}

	resp, err := h.svc.GetMetric(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metricResponsePayload(resp))
}

func (h *Handler) ListMetrics(c *gin.Context) {
	metricsMap := h.svc.ListAll()

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)

	for _, g := range metricsMap["gauge"] {
		_, _ = c.Writer.Write([]byte(g + "<br>\n"))
	}
	for _, cM := range metricsMap["counter"] {
		_, _ = c.Writer.Write([]byte(cM + "<br>\n"))
	}
}

func (h *Handler) NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
}

func (h *Handler) BadRequestHandler(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
}

func (h *Handler) updateFromPath(c *gin.Context, mtype, name, valueStr string) (gin.H, error) {
	switch mtype {
	case "gauge":
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value for gauge"})
			return nil, err
		}
		var resp metrics.MetricResponse
		resp, err = h.svc.UpdateGauge(name, value)
		return metricResponsePayload(resp), err
	case "counter":
		delta, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value for counter"})
			return nil, err
		}
		var resp metrics.MetricResponse
		resp, err = h.svc.UpdateCounter(name, delta)
		return metricResponsePayload(resp), err
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid metric type"})
		return nil, nil
	}
}
