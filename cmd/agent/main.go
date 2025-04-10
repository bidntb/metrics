package main

import (
	"bidntb/metrics/internal/collector"
	"bidntb/metrics/internal/nconfig"

	"fmt"
	"time"
)

func main() {
	address, reportInterval, pollInterval := nconfig.ParseFlags()

	metrics := collector.NewMetrics()
	serverURL := fmt.Sprintf("http://%s", address)

	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)

	fmt.Printf("Starting metrics collector:\n")
	fmt.Printf("Server URL: %s\n", serverURL)
	fmt.Printf("Poll interval: %d seconds\n", pollInterval)
	fmt.Printf("Report interval: %d seconds\n", reportInterval)

	for {
		select {
		case <-pollTicker.C:
			metrics.UpdateMetrics()
		case <-reportTicker.C:
			if err := metrics.SendMetrics(serverURL); err != nil {
				fmt.Printf("Error sending metrics: %v\n", err)
			}
		}
	}
}
