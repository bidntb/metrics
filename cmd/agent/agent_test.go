package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	metrics := NewMetrics()

	if metrics == nil {
		t.Error("NewMetrics() returned nil")
	}

	if metrics.gauges == nil {
		t.Error("gauges map was not initialized")
	}

	if metrics.counters == nil {
		t.Error("counters map was not initialized")
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
