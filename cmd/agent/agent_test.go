package main

import (
	"bidntb/metrics/internal/agent/collector"
	"bidntb/metrics/internal/agent/reporter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	metrics := collector.NewMetrics()

	if metrics == nil {
		t.Fatal("NewMetrics() returned nil")
	}

	if metrics.Gauges == nil {
		t.Error("gauges map was not initialized")
	}

	if metrics.Counters == nil {
		t.Error("counters map was not initialized")
	}
}

func TestUpdateMetrics(t *testing.T) {
	metrics := collector.NewMetrics()
	metrics.UpdateMetrics()

	expectedGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for _, gauge := range expectedGauges {
		if _, exists := metrics.Gauges[gauge]; !exists {
			t.Errorf("Expected gauge %s not found", gauge)
		}
	}

	if _, exists := metrics.Counters["PollCount"]; !exists {
		t.Error("PollCount counter not found")
	}

	if metrics.Counters["PollCount"] != 1 {
		t.Errorf("Expected PollCount to be 1, got %d", metrics.Counters["PollCount"])
	}
}

func TestSendMetrics(t *testing.T) {
	metrics := collector.NewMetrics()

	metrics.Gauges["TestGauge"] = 123.45
	metrics.Counters["TestCounter"] = 42

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := reporter.SendMetricsJSON(server.URL, metrics)
	if err != nil {
		t.Errorf("sendMetrics() returned unexpected error: %v", err)
	}
}

func TestSendMetricsError(t *testing.T) {
	metrics := collector.NewMetrics()
	metrics.Gauges["TestGauge"] = 123.45

	err := reporter.SendMetricsJSON("http://localhost:12345", metrics)
	if err == nil {
		t.Error("Expected error when sending to non-existent server, got nil")
	}
}

func TestMultipleUpdateMetrics(t *testing.T) {
	metrics := collector.NewMetrics()

	for i := 0; i < 3; i++ {
		metrics.UpdateMetrics()
	}

	if metrics.Counters["PollCount"] != 3 {
		t.Errorf("Expected PollCount to be 3, got %d", metrics.Counters["PollCount"])
	}
}

func TestRandomValueChanges(t *testing.T) {
	metrics := collector.NewMetrics()

	metrics.UpdateMetrics()
	firstValue := metrics.Gauges["RandomValue"]

	metrics.UpdateMetrics()
	secondValue := metrics.Gauges["RandomValue"]

	if firstValue == secondValue {
		t.Error("RandomValue should change between updates")
	}
}
