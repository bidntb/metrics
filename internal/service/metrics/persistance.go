package metrics

import (
	"encoding/json"
	"fmt"
	"os"

	"bidntb/metrics/internal/storage"
)

type Snapshot struct {
	Gauges   []storage.GaugeMetric   `json:"gauges"`
	Counters []storage.CounterMetric `json:"counters"`
}

func (s *Service) SaveTo(path string) error {
	snap := Snapshot{
		Gauges:   s.storage.GetAllGauge(),
		Counters: s.storage.GetAllCounter(),
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write tmp snapshot: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename snapshot: %w", err)
	}
	return nil
}

func (s *Service) LoadFrom(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read snapshot: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("unmarshal snapshot: %w", err)
	}

	for _, g := range snap.Gauges {
		s.storage.AddGauge(g)
	}
	for _, c := range snap.Counters {
		s.storage.AddCounter(c)
	}

	return nil
}
