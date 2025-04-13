package main

import (
	"bidntb/metrics/internal/agent/collector"
	"bidntb/metrics/internal/agent/reporter"
	"bidntb/metrics/internal/nconfig"

	"fmt"
	"time"
)

func MetricService() {
	serverURL, reportInterval, pollInterval := nconfig.ParseFlags()

	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	metrics := collector.NewMetrics()

	for {
		select {
		case <-pollTicker.C:
			metrics.UpdateMetrics()
		case <-reportTicker.C:
			if err := reporter.SendMetrics(serverURL, metrics); err != nil {
				fmt.Printf("Error sending metrics: %v\n", err)
			}
		}
	}

}

func main() {
	MetricService()
}
