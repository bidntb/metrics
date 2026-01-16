package metrics

import (
	"fmt"
	"time"

	"bidntb/metrics/internal/storage"
)

type Service struct {
	storage storage.Interface
}

func NewService(storage storage.Interface) *Service {
	return &Service{storage: storage}
}

type UpdateMetricRequest struct {
	Type  string  `json:"type" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Value float64 `json:"value" binding:"required"`
}

type UpdateCounterRequest struct {
	Delta int64 `json:"value" binding:"required"`
}

type UpdateGaugeRequest struct {
	Value float64 `json:"value" binding:"required"`
}

type GetMetricRequest struct {
	ID    string `json:"id" binding:"required"`
	MType string `json:"MType" binding:"required,oneof=gauge counter"`
}

type MetricResponse struct {
	ID    string  `json:"id"`
	MType string  `json:"MType"`
	Delta *int64  `json:"delta,omitempty"`
	Value float64 `json:"value"`
}

func (s *Service) UpdateMetric(req UpdateMetricRequest) (*MetricResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("missing name")
	}

	switch req.Type {
	case "gauge":
		return s.UpdateGauge(req.Name, req.Value)
	case "counter":
		return s.UpdateCounter(req.Name, int64(req.Value))
	default:
		return nil, fmt.Errorf("invalid metric type: %s", req.Type)
	}
}

func (s *Service) GetMetric(req GetMetricRequest) (*MetricResponse, bool) {
	switch req.MType {
	case "gauge":
		return s.getGaugeByID(req.ID)
	case "counter":
		return s.getCounterByID(req.ID)
	}
	return nil, false
}

func (s *Service) GetMetricValue(mtype, name string) (string, bool) {
	switch mtype {
	case "gauge":
		return s.getGaugeValue(name)
	case "counter":
		return s.getCounterValue(name)
	}
	return "", false
}

func (s *Service) UpdateGauge(name string, value float64) (*MetricResponse, error) {
	metric := storage.GaugeMetric{
		ID:         int(time.Now().UnixNano()),
		MetricName: name,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}
	s.storage.AddGaugeMetric(metric)

	if m, found := s.storage.GetGaugeMetric(name); found {
		return &MetricResponse{
			ID:    fmt.Sprintf("%v", m.ID),
			MType: "gauge",
			Value: m.Value,
		}, nil
	}
	return nil, fmt.Errorf("gauge metric not found after update")
}

func (s *Service) UpdateCounter(name string, delta int64) (*MetricResponse, error) {
	last, exists := s.storage.GetCounterMetric(name)
	value := delta
	if exists {
		value = last.Value + delta
	}

	metric := storage.CounterMetric{
		ID:         int(time.Now().UnixNano()),
		MetricName: name,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}
	s.storage.AddCounterMetric(metric)

	if m, found := s.storage.GetCounterMetric(name); found {
		deltaVal := m.Value
		return &MetricResponse{
			ID:    fmt.Sprintf("%v", m.ID),
			MType: "counter",
			Delta: &deltaVal,
			Value: float64(m.Value),
		}, nil
	}
	return nil, fmt.Errorf("counter metric not found after update")
}

func (s *Service) getGaugeByID(id string) (*MetricResponse, bool) {
	if metric, found := s.storage.GetGaugeMetricByID(id); found {
		return &MetricResponse{
			ID:    id,
			MType: "gauge",
			Value: metric.Value,
		}, true
	}
	return nil, false
}

func (s *Service) getCounterByID(id string) (*MetricResponse, bool) {
	if metric, found := s.storage.GetCounterMetricByID(id); found {
		delta := metric.Value
		return &MetricResponse{
			ID:    id,
			MType: "counter",
			Delta: &delta,
			Value: float64(metric.Value),
		}, true
	}
	return nil, false
}

func (s *Service) getGaugeValue(name string) (string, bool) {
	if m, ok := s.storage.GetGaugeMetric(name); ok {
		return fmt.Sprintf("%.3f", m.Value), true
	}
	return "", false
}

func (s *Service) getCounterValue(name string) (string, bool) {
	if m, ok := s.storage.GetCounterMetric(name); ok {
		return fmt.Sprintf("%d", m.Value), true
	}
	return "", false
}

func (s *Service) ListAll() map[string][]string {
	gauges := make([]string, 0)
	for _, m := range s.storage.GetGaugeMetrics() {
		gauges = append(gauges, fmt.Sprintf("%d: %s - %.3f", m.ID, m.MetricName, m.Value))
	}

	counters := make([]string, 0)
	for _, m := range s.storage.GetCounterMetrics() {
		counters = append(counters, fmt.Sprintf("%d: %s - %d", m.ID, m.MetricName, m.Value))
	}

	return map[string][]string{
		"gauge":   gauges,
		"counter": counters,
	}
}
