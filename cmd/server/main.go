package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GaugeMetric struct {
	id         int
	timestamp  int64
	metricName string
	value      float64
}

type CounterMetric struct {
	id         int
	metricName string
	timestamp  int64
	value      int64
}

type MemStorage struct {
	GaugeMetrics   []GaugeMetric
	CounterMetrics []CounterMetric
}

func (s *MemStorage) AddGaugeMetric(metric GaugeMetric) {
	s.GaugeMetrics = append(s.GaugeMetrics, metric)
}

func (s *MemStorage) AddCounterMetric(metric CounterMetric) {
	s.CounterMetrics = append(s.CounterMetrics, metric)
}

var storage = &MemStorage{
	GaugeMetrics:   make([]GaugeMetric, 0),
	CounterMetrics: make([]CounterMetric, 0),
}

var counterMap = make(map[string]int64)

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

	gaugeMetric := GaugeMetric{
		id:         len(storage.GaugeMetrics),
		metricName: metricName,
		timestamp:  time.Now().Unix(),
		value:      value,
	}

	storage.AddGaugeMetric(gaugeMetric)
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

	lastValue, exists := counterMap[metricName]
	if exists {
		value = value + lastValue
	}

	counterMetric := CounterMetric{
		id:         len(storage.CounterMetrics),
		metricName: metricName,
		timestamp:  time.Now().Unix(),
		value:      value,
	}

	storage.AddCounterMetric(counterMetric)
	counterMap[metricName] = value
	c.Status(http.StatusOK)
}

func main() {
	router := gin.Default()

	// Routes
	router.POST("/update/gauge/:name/:value", gaugeHandler)
	router.POST("/update/counter/:name/:value", counterHandler)

	// Catch-all route for unknown paths
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
	})

	// Start server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
