package reporter

import (
	"bidntb/metrics/internal/agent/collector"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		resp.Body.Close()
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
		resp.Body.Close()
	}
	return nil
}

func SendMetricsJSON(serverURL string, metrics *Metrics) error {
	client := &http.Client{}

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/update/", serverURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
