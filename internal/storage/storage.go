package storage

import "sync"

type Interface interface {
	AddGaugeMetric(metric GaugeMetric)
	AddCounterMetric(metric CounterMetric)
	GetLastGauge(name string) (*GaugeMetric, bool)
	GetLastCounter(name string) (*CounterMetric, bool)
	GetAllGaugeMetrics() []GaugeMetric
	GetAllCounterMetrics() []CounterMetric
}

type Metric interface {
	GetName() string
}

type GaugeMetric struct {
	ID         int     `json:"id"`
	MetricName string  `json:"metric_name"`
	Timestamp  int64   `json:"timestamp"`
	Value      float64 `json:"value"`
}

func (m GaugeMetric) GetName() string { return m.MetricName }

type CounterMetric struct {
	ID         int    `json:"id"`
	MetricName string `json:"metric_name"`
	Timestamp  int64  `json:"timestamp"`
	Value      int64  `json:"value"`
}

func (m CounterMetric) GetName() string { return m.MetricName }

type MemStorage struct {
	gaugeMu   sync.RWMutex
	counterMu sync.RWMutex

	gauges   map[string][]GaugeMetric
	counters map[string][]CounterMetric
}

func NewMemStorage() Interface {
	return &MemStorage{
		gauges:   make(map[string][]GaugeMetric),
		counters: make(map[string][]CounterMetric),
	}
}

func (s *MemStorage) AddGaugeMetric(metric GaugeMetric) {
	s.gaugeMu.Lock()
	defer s.gaugeMu.Unlock()

	s.gauges[metric.MetricName] = append(s.gauges[metric.MetricName], metric)
}

func (s *MemStorage) AddCounterMetric(metric CounterMetric) {
	s.counterMu.Lock()
	defer s.counterMu.Unlock()

	s.counters[metric.MetricName] = append(s.counters[metric.MetricName], metric)
}

func (s *MemStorage) GetLastGauge(name string) (*GaugeMetric, bool) {
	s.gaugeMu.RLock()
	defer s.gaugeMu.RUnlock()

	metrics, ok := s.gauges[name]
	if !ok || len(metrics) == 0 {
		return nil, false
	}
	last := metrics[len(metrics)-1]
	return &last, true
}

func (s *MemStorage) GetLastCounter(name string) (*CounterMetric, bool) {
	s.counterMu.RLock()
	defer s.counterMu.RUnlock()

	metrics, ok := s.counters[name]
	if !ok || len(metrics) == 0 {
		return nil, false
	}
	last := metrics[len(metrics)-1]
	return &last, true
}

func (s *MemStorage) GetAllGaugeMetrics() []GaugeMetric {
	s.gaugeMu.RLock()
	defer s.gaugeMu.RUnlock()

	var all []GaugeMetric
	for _, metrics := range s.gauges {
		all = append(all, metrics...)
	}
	return all
}

func (s *MemStorage) GetAllCounterMetrics() []CounterMetric {
	s.counterMu.RLock()
	defer s.counterMu.RUnlock()

	var all []CounterMetric
	for _, metrics := range s.counters {
		all = append(all, metrics...)
	}
	return all
}
