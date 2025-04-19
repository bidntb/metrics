package storage

type StorageInterface interface {
	AddGaugeMetric(metric GaugeMetric)
	AddCounterMetric(metric CounterMetric)
	GetGaugeMetric(name string) (*GaugeMetric, bool)
	GetCounterMetric(name string) (*CounterMetric, bool)
	GetGaugeMetrics() []GaugeMetric
	GetCounterMetrics() []CounterMetric
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

func NewMemStorage() StorageInterface {
	return &MemStorage{
		GaugeMetrics:   make([]GaugeMetric, 0),
		CounterMetrics: make([]CounterMetric, 0),
	}
}
