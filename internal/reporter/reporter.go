package reporter

import (
	"fmt"
	"net/http"
)

type Metrics struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *Metrics) SendMetrics(serverURL string) error {
	client := &http.Client{}

	for name, value := range m.gauges {
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

	for name, value := range m.counters {
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
