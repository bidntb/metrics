package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"bidntb/metrics/internal/nconfig"
	"bidntb/metrics/internal/storage"
)

func gaugeHandler(c *gin.Context) {
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

func counterHandler(c *gin.Context) {
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

func getMetricHandler(c *gin.Context) {
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

func indexHandler(c *gin.Context) {
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

func main() {

	ServerAddress := nconfig.GetServerAddress()

	router := gin.Default()

	router.GET("/", indexHandler)
	router.GET("/value/:type/:name", getMetricHandler)

	router.POST("/update/counter/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/counter", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/gauge/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/gauge", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
	router.POST("/update/:wrong", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	})
	router.POST("/update/:wrong/*any", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	})

	router.POST("/update/gauge/:name/:value", gaugeHandler)
	router.POST("/update/counter/:name/:value", counterHandler)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	if err := router.Run(ServerAddress); err != nil {
		panic(err)
	}
}
