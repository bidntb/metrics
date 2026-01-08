package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/storage"
)

type Handler struct {
	storage storage.StorageInterface
}

func NewHandler(storage storage.StorageInterface) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) AddGaugeHandler(c *gin.Context) {
	metricName := c.Param("name")
	valueStr := c.Param("value")

	if metricName == "" || valueStr == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid URL"})
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Value"})
		return
	}

	gaugeMetric := storage.GaugeMetric{
		ID:         len(h.storage.GetGaugeMetrics()),
		MetricName: metricName,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}

	h.storage.AddGaugeMetric(gaugeMetric)
	c.Status(http.StatusOK)
}

func (h *Handler) AddCounterHandler(c *gin.Context) {
	metricName := c.Param("name")
	valueStr := c.Param("value")

	if metricName == "" || valueStr == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid URL"})
		return
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Value"})
		return
	}

	lastMetric, exists := h.storage.GetCounterMetric(metricName)
	if exists {
		value = value + lastMetric.Value
	}

	counterMetric := storage.CounterMetric{
		ID:         len(h.storage.GetCounterMetrics()),
		MetricName: metricName,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}

	h.storage.AddCounterMetric(counterMetric)
	c.Status(http.StatusOK)
}

func (h *Handler) ValueHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	switch metricType {
	case "gauge":
		if metric, found := h.storage.GetGaugeMetric(metricName); found {
			c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
			return
		}
	case "counter":
		if metric, found := h.storage.GetCounterMetric(metricName); found {
			c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func (h *Handler) IndexHandler(c *gin.Context) {
	var gaugeMetrics []string
	for _, metric := range h.storage.GetGaugeMetrics() {
		gaugeMetrics = append(gaugeMetrics, fmt.Sprintf("%d: %s - %.2f", metric.ID, metric.MetricName, metric.Value))
	}

	var counterMetrics []string
	for _, metric := range h.storage.GetCounterMetrics() {
		counterMetrics = append(counterMetrics, fmt.Sprintf("%d: %s - %d", metric.ID, metric.MetricName, metric.Value))
	}

	c.JSON(http.StatusOK, gin.H{
		"gauge_metrics":   gaugeMetrics,
		"counter_metrics": counterMetrics,
	})
}
