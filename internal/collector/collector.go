package collector

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
)

type Metrics struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *Metrics) UpdateMetrics() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	m.Gauges["Alloc"] = float64(stats.Alloc)
	m.Gauges["BuckHashSys"] = float64(stats.BuckHashSys)
	m.Gauges["Frees"] = float64(stats.Frees)
	m.Gauges["GCCPUFraction"] = stats.GCCPUFraction
	m.Gauges["GCSys"] = float64(stats.GCSys)
	m.Gauges["HeapAlloc"] = float64(stats.HeapAlloc)
	m.Gauges["HeapIdle"] = float64(stats.HeapIdle)
	m.Gauges["HeapInuse"] = float64(stats.HeapInuse)
	m.Gauges["HeapObjects"] = float64(stats.HeapObjects)
	m.Gauges["HeapReleased"] = float64(stats.HeapReleased)
	m.Gauges["HeapSys"] = float64(stats.HeapSys)
	m.Gauges["LastGC"] = float64(stats.LastGC)
	m.Gauges["Lookups"] = float64(stats.Lookups)
	m.Gauges["MCacheInuse"] = float64(stats.MCacheInuse)
	m.Gauges["MCacheSys"] = float64(stats.MCacheSys)
	m.Gauges["MSpanInuse"] = float64(stats.MSpanInuse)
	m.Gauges["MSpanSys"] = float64(stats.MSpanSys)
	m.Gauges["Mallocs"] = float64(stats.Mallocs)
	m.Gauges["NextGC"] = float64(stats.NextGC)
	m.Gauges["NumForcedGC"] = float64(stats.NumForcedGC)
	m.Gauges["NumGC"] = float64(stats.NumGC)
	m.Gauges["OtherSys"] = float64(stats.OtherSys)
	m.Gauges["PauseTotalNs"] = float64(stats.PauseTotalNs)
	m.Gauges["StackInuse"] = float64(stats.StackInuse)
	m.Gauges["StackSys"] = float64(stats.StackSys)
	m.Gauges["Sys"] = float64(stats.Sys)
	m.Gauges["TotalAlloc"] = float64(stats.TotalAlloc)

	m.Gauges["RandomValue"] = rand.Float64()

	m.Counters["PollCount"]++
}

func (m *Metrics) SendMetrics(serverURL string) error {
	client := &http.Client{}

	for name, value := range m.Gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", serverURL, name, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	for name, value := range m.Counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d", serverURL, name, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}
	return nil
}
