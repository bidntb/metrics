package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	metrics := NewMetrics()

	if metrics == nil {
		t.Fatal("NewMetrics() returned nil")
	}

	if metrics.gauges == nil {
		t.Error("gauges map was not initialized")
	}

	if metrics.counters == nil {
		t.Error("counters map was not initialized")
	}
}

func TestUpdateMetrics(t *testing.T) {
	metrics := NewMetrics()
	metrics.updateMetrics()

	expectedGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for _, gauge := range expectedGauges {
		if _, exists := metrics.gauges[gauge]; !exists {
			t.Errorf("Expected gauge %s not found", gauge)
		}
	}

	if _, exists := metrics.counters["PollCount"]; !exists {
		t.Error("PollCount counter not found")
	}

	if metrics.counters["PollCount"] != 1 {
		t.Errorf("Expected PollCount to be 1, got %d", metrics.counters["PollCount"])
	}
}

func TestSendMetrics(t *testing.T) {
	metrics := NewMetrics()
	metrics.gauges["TestGauge"] = 123.45
	metrics.counters["TestCounter"] = 42

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := metrics.sendMetrics(server.URL)
	if err != nil {
		t.Errorf("sendMetrics() returned unexpected error: %v", err)
	}
}

func TestSendMetricsError(t *testing.T) {
	metrics := NewMetrics()
	metrics.gauges["TestGauge"] = 123.45

	err := metrics.sendMetrics("http://localhost:12345")
	if err == nil {
		t.Error("Expected error when sending to non-existent server, got nil")
	}
}

func TestMultipleUpdateMetrics(t *testing.T) {
	metrics := NewMetrics()

	for i := 0; i < 3; i++ {
		metrics.updateMetrics()
	}

	if metrics.counters["PollCount"] != 3 {
		t.Errorf("Expected PollCount to be 3, got %d", metrics.counters["PollCount"])
	}
}

func TestRandomValueChanges(t *testing.T) {
	metrics := NewMetrics()

	metrics.updateMetrics()
	firstValue := metrics.gauges["RandomValue"]

	metrics.updateMetrics()
	secondValue := metrics.gauges["RandomValue"]

	if firstValue == secondValue {
		t.Error("RandomValue should change between updates")
	}
}
