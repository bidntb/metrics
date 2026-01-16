package storage

import "fmt"

type Interface interface {
	AddGaugeMetric(metric GaugeMetric)
	AddCounterMetric(metric CounterMetric)
	GetGaugeMetric(name string) (*GaugeMetric, bool)
	GetCounterMetric(name string) (*CounterMetric, bool)
	GetGaugeMetrics() []GaugeMetric
	GetCounterMetrics() []CounterMetric
	GetGaugeMetricByID(id string) (*GaugeMetric, bool)
	GetCounterMetricByID(id string) (*CounterMetric, bool)
}

type GaugeMetric struct {
	ID         int
	Timestamp  int64
	MetricName string
	Value      float64
}

type CounterMetric struct {
	ID         int
	MetricName string
	Timestamp  int64
	Value      int64
	LastValue  int64
}

type MemStorage struct {
	GaugeMetrics   []GaugeMetric
	CounterMetrics []CounterMetric
}

func (s *MemStorage) GetGaugeMetricByID(id string) (*GaugeMetric, bool) {
	for _, metric := range s.GaugeMetrics {
		if fmt.Sprintf("%v", metric.ID) == id {
			return &metric, true
		}
	}
	return nil, false
}

func (s *MemStorage) GetCounterMetricByID(id string) (*CounterMetric, bool) {
	for _, metric := range s.CounterMetrics {
		if fmt.Sprintf("%v", metric.ID) == id {
			return &metric, true
		}
	}
	return nil, false
}

func (s *MemStorage) AddGaugeMetric(metric GaugeMetric) {
	s.GaugeMetrics = append(s.GaugeMetrics, metric)
}

func (s *MemStorage) AddCounterMetric(metric CounterMetric) {
	s.CounterMetrics = append(s.CounterMetrics, metric)

}

func (s *MemStorage) GetGaugeMetric(name string) (*GaugeMetric, bool) {
	for i := len(s.GaugeMetrics) - 1; i >= 0; i-- {
		if s.GaugeMetrics[i].MetricName == name {
			return &s.GaugeMetrics[i], true
		}
	}
	return nil, false
}

func (s *MemStorage) GetCounterMetric(name string) (*CounterMetric, bool) {
	for i := len(s.CounterMetrics) - 1; i >= 0; i-- {
		if s.CounterMetrics[i].MetricName == name {
			return &s.CounterMetrics[i], true
		}
	}
	return nil, false
}

func (s *MemStorage) GetGaugeMetrics() []GaugeMetric {
	return s.GaugeMetrics
}

func (s *MemStorage) GetCounterMetrics() []CounterMetric {
	return s.CounterMetrics
}

func NewMemStorage() Interface {
	return &MemStorage{
		GaugeMetrics:   make([]GaugeMetric, 0),
		CounterMetrics: make([]CounterMetric, 0),
	}
}
