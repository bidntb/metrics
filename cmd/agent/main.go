package main

import (
	"bidntb/metrics/internal/agent/collector"
	"bidntb/metrics/internal/agent/reporter"
	"bidntb/metrics/internal/nconfig"

	"fmt"
	"time"
)

func MetricService() {
	cfg := nconfig.ParseConfig()

	serverURL := fmt.Sprintf("https://%s", cfg.ServerAddr)
	reportInterval := cfg.ReportInterval
	pollInterval := cfg.PollInterval
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	metrics := collector.NewMetrics()

	for {
		select {
		case <-pollTicker.C:
			metrics.UpdateMetrics()
		case <-reportTicker.C:
			if err := reporter.SendMetricsJSON(serverURL, metrics); err != nil {
				fmt.Printf("Error sending metrics through path: %v\n", err)
				err := reporter.SendMetrics(serverURL, metrics)
				if err != nil {
					fmt.Printf("Error sending metrics through JSON: %v\n", err)
				}
			}
		}
	}

}

func main() {
	MetricService()
}
