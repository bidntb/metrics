package main

import (
	"fmt"
	"net/http"
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
var gaugeMap = make(map[string]float64)

func gaugeHandler(c *gin.Context) {
	var requestBody struct {
		MetricName string  `json:"metric_name"`
		Value      float64 `json:"value"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	idStr := len(storage.GaugeMetrics)
	gaugeMetric := GaugeMetric{
		id:         idStr,
		metricName: requestBody.MetricName,
		timestamp:  time.Now().Unix(),
		value:      requestBody.Value,
	}

	storage.AddGaugeMetric(gaugeMetric)
	c.JSON(http.StatusOK, gin.H{"id": idStr})
}

func counterHandler(c *gin.Context) {
	var requestBody struct {
		MetricName string `json:"metric_name"`
		Value      int64  `json:"value"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	value, check := counterMap[requestBody.MetricName]
	if check {
		requestBody.Value += value
	}

	idStr := len(storage.CounterMetrics)
	counterMetric := CounterMetric{
		id:         idStr,
		metricName: requestBody.MetricName,
		timestamp:  time.Now().Unix(),
		value:      requestBody.Value,
	}

	storage.AddCounterMetric(counterMetric)
	counterMap[requestBody.MetricName] = requestBody.Value
	c.JSON(http.StatusOK, gin.H{"id": idStr})
}

func valueHandler(c *gin.Context) {
	metricType := c.Param("metric_type")
	metricName := c.Param("metric_name")

	var totalValue float64
	if metricType == "gauge" {
		value, check := gaugeMap[metricName]
		if check {
			totalValue = value
		}
	} else if metricType == "counter" {
		value, check := counterMap[metricName]
		if check {
			totalValue = float64(value)
		}
	}

	c.JSON(http.StatusOK, gin.H{"value": totalValue})
}

func MainHandler(c *gin.Context) {
	var gaugeMetrics []string
	for _, metric := range storage.GaugeMetrics {
		gaugeMetrics = append(gaugeMetrics, fmt.Sprintf("%d: %s - %.2f", metric.id, metric.metricName, metric.value))
	}

	var counterMetrics []string
	for _, metric := range storage.CounterMetrics {
		counterMetrics = append(counterMetrics, fmt.Sprintf("%d: %s - %d", metric.id, metric.metricName, metric.value))
	}

	c.JSON(http.StatusOK, gin.H{
		"gauge_metrics":   gaugeMetrics,
		"counter_metrics": counterMetrics,
	})
}

func main() {
	r := gin.Default()

	r.POST("/update/gauge/", gaugeHandler)
	r.POST("/update/counter/", counterHandler)
	r.GET("/value/:metric_type/:metric_name", valueHandler)
	r.GET("/", MainHandler)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
