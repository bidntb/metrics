package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Metrics struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *Metrics) updateMetrics() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	m.gauges["Alloc"] = float64(stats.Alloc)
	m.gauges["BuckHashSys"] = float64(stats.BuckHashSys)
	m.gauges["Frees"] = float64(stats.Frees)
	m.gauges["GCCPUFraction"] = stats.GCCPUFraction
	m.gauges["GCSys"] = float64(stats.GCSys)
	m.gauges["HeapAlloc"] = float64(stats.HeapAlloc)
	m.gauges["HeapIdle"] = float64(stats.HeapIdle)
	m.gauges["HeapInuse"] = float64(stats.HeapInuse)
	m.gauges["HeapObjects"] = float64(stats.HeapObjects)
	m.gauges["HeapReleased"] = float64(stats.HeapReleased)
	m.gauges["HeapSys"] = float64(stats.HeapSys)
	m.gauges["LastGC"] = float64(stats.LastGC)
	m.gauges["Lookups"] = float64(stats.Lookups)
	m.gauges["MCacheInuse"] = float64(stats.MCacheInuse)
	m.gauges["MCacheSys"] = float64(stats.MCacheSys)
	m.gauges["MSpanInuse"] = float64(stats.MSpanInuse)
	m.gauges["MSpanSys"] = float64(stats.MSpanSys)
	m.gauges["Mallocs"] = float64(stats.Mallocs)
	m.gauges["NextGC"] = float64(stats.NextGC)
	m.gauges["NumForcedGC"] = float64(stats.NumForcedGC)
	m.gauges["NumGC"] = float64(stats.NumGC)
	m.gauges["OtherSys"] = float64(stats.OtherSys)
	m.gauges["PauseTotalNs"] = float64(stats.PauseTotalNs)
	m.gauges["StackInuse"] = float64(stats.StackInuse)
	m.gauges["StackSys"] = float64(stats.StackSys)
	m.gauges["Sys"] = float64(stats.Sys)
	m.gauges["TotalAlloc"] = float64(stats.TotalAlloc)

	m.gauges["RandomValue"] = rand.Float64()

	m.counters["PollCount"]++
}

func (m *Metrics) sendMetrics(serverURL string) error {
	client := &http.Client{}

	for name, value := range m.gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", serverURL, name, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
		req.Body.Close()
	}

	for name, value := range m.counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d", serverURL, name, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
		return nil
		req.Body.Close()
	}
	return nil
}

func main() {
	metrics := NewMetrics()
	serverURL := "http://localhost:8080"

	pollTicker := time.NewTicker(2 * time.Second)
	reportTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-pollTicker.C:
			metrics.updateMetrics()
		case <-reportTicker.C:
			if err := metrics.sendMetrics(serverURL); err != nil {
				fmt.Printf("Error sending metrics: %v\n", err)
			}
		}
	}
}
