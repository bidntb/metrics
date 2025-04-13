package reporter

import (
	"bidntb/metrics/internal/agent/collector"

	"fmt"
	"net/http"
)

type Metrics = collector.Metrics

func SendMetrics(serverURL string, metrics *Metrics) error {
	client := &http.Client{}

	for name, value := range metrics.Gauges {
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

	for name, value := range metrics.Counters {
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
