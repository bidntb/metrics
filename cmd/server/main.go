package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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

func gaugeHandler(res http.ResponseWriter, req *http.Request) {
	storage := MemStorage{}

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/update/gauge/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	metricName := parts[0]
	idStr := len(storage.GaugeMetrics)
	value, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		http.Error(res, "Invalid Value", http.StatusBadRequest)
		return
	}

	gaugeMetric := GaugeMetric{
		id:         idStr,
		metricName: metricName,
		timestamp:  time.Now().Unix(),
		value:      value,
	}

	storage.AddGaugeMetric(gaugeMetric)
	res.WriteHeader(http.StatusOK)
}

func counterHandler(res http.ResponseWriter, req *http.Request) {
	storage := MemStorage{}

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/update/counter/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	metricName := parts[0]
	idStr := len(storage.CounterMetrics)
	value, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		http.Error(res, "Invalid Value", http.StatusBadRequest)
		return
	}
	for _, metric := range storage.CounterMetrics {
		if metric.metricName == metricName {
			value = value + metric.value
		}
	}

	CounterMetric := CounterMetric{
		id:         idStr,
		metricName: metricName,
		timestamp:  time.Now().Unix(),
		value:      value,
	}

	storage.AddCounterMetric(CounterMetric)
	res.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, gaugeHandler)
	mux.HandleFunc(`/update/counter/`, counterHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
