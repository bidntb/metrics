package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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

func (s *MemStorage) GetGaugeMetric(name string) (*GaugeMetric, bool) {
	for i := len(s.GaugeMetrics) - 1; i >= 0; i-- {
		if s.GaugeMetrics[i].metricName == name {
			return &s.GaugeMetrics[i], true
		}
	}
	return nil, false
}

func (s *MemStorage) GetCounterMetric(name string) (*CounterMetric, bool) {
	for i := len(s.CounterMetrics) - 1; i >= 0; i-- {
		if s.CounterMetrics[i].metricName == name {
			return &s.CounterMetrics[i], true
		}
	}
	return nil, false
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

func getMetricHandler(c *gin.Context) {
	metricType := c.Param("type")
	metricName := c.Param("name")

	switch metricType {
	case "gauge":
		if metric, found := storage.GetGaugeMetric(metricName); found {
			c.String(http.StatusOK, fmt.Sprintf("%v", metric.value))
			return
		}
	case "counter":
		if metric, found := storage.GetCounterMetric(metricName); found {
			c.String(http.StatusOK, fmt.Sprintf("%v", metric.value))
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func indexHandler(c *gin.Context) {
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

func NoRedirect(c *gin.Context) {
	c.Request.Header.Set("Connection", "close")
	c.Next()
}

func nonRegisteredPathHandler(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
}

func main() {
	defaultAddress := "localhost:8080"

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		defaultAddress = envAddress
	}

	serverAddress := flag.String("a", defaultAddress, "HTTP server endpoint address")
	flag.Parse()

	finalAddress := *serverAddress
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		finalAddress = envAddress
	}

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

	if err := router.Run(finalAddress); err != nil {
		panic(err)
	}
}
