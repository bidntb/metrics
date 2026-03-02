package reporter

import (
	"bidntb/metrics/internal/agent/collector"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Metrics = collector.Metrics

func SendMetrics(serverURL string, metrics *Metrics) error {
	client := &http.Client{}

	for name, value := range metrics.Gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%f", serverURL, name, value)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		er := resp.Body.Close()
		if er != nil {
			return er
		}
	}

	for name, value := range metrics.Counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d", serverURL, name, value)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		er := resp.Body.Close()
		if er != nil {
			return er
		}
	}
	return nil
}

type updatePayload struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

const sendRetries = 5
const sendRetryDelay = 200 * time.Millisecond

func SendMetricsJSON(serverURL string, metrics *Metrics) error {
	client := &http.Client{Timeout: 5 * time.Second}
	baseURL := fmt.Sprintf("%s/update/", serverURL)

	sendOne := func(data []byte) error {
		var lastErr error
		for attempt := 0; attempt < sendRetries; attempt++ {
			req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(data))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				lastErr = err
				time.Sleep(sendRetryDelay)
				continue
			}
			er := resp.Body.Close()
			if er != nil {
				return er
			}
			return nil
		}
		return lastErr
	}

	for name, value := range metrics.Gauges {
		body := updatePayload{ID: name, MType: "gauge", Value: &value}
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		if err := sendOne(data); err != nil {
			return err
		}
	}

	for name, value := range metrics.Counters {
		body := updatePayload{ID: name, MType: "counter", Delta: &value}
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		if err := sendOne(data); err != nil {
			return err
		}
	}

	return nil
}
