package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/storage"
)

func AddGaugeHandler(c *gin.Context) {
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
		ID:         len(storage.Storage.GaugeMetrics),
		MetricName: metricName,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}

	storage.Storage.AddGaugeMetric(gaugeMetric)
	c.Status(http.StatusOK)
}

func AddCounterHandler(c *gin.Context) {
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

	lastValue, exists := storage.CounterMap[metricName]
	if exists {
		value = value + lastValue
	}

	counterMetric := storage.CounterMetric{
		ID:         len(storage.Storage.CounterMetrics),
		MetricName: metricName,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}

	storage.Storage.AddCounterMetric(counterMetric)
	storage.CounterMap[metricName] = value
	c.Status(http.StatusOK)
}

func ValueHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	switch metricType {
	case "gauge":
		if metric, found := storage.Storage.GetGaugeMetric(metricName); found {
			c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
			return
		}
	case "counter":
		if metric, found := storage.Storage.GetCounterMetric(metricName); found {
			c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func IndexHandler(c *gin.Context) {
	var gaugeMetrics []string
	for _, metric := range storage.Storage.GaugeMetrics {
		gaugeMetrics = append(gaugeMetrics, fmt.Sprintf("%d: %s - %.2f", metric.ID, metric.MetricName, metric.Value))
	}

	var counterMetrics []string
	for _, metric := range storage.Storage.CounterMetrics {
		counterMetrics = append(counterMetrics, fmt.Sprintf("%d: %s - %d", metric.ID, metric.MetricName, metric.Value))
	}

	c.JSON(http.StatusOK, gin.H{
		"gauge_metrics":   gaugeMetrics,
		"counter_metrics": counterMetrics,
	})
}
