package metrics

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bidntb/metrics/internal/storage"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Service struct {
	storage storage.Interface
}

func NewService(storage storage.Interface) *Service {
	return &Service{storage: storage}
}

type UpdateMetricRequest struct {
	ID    string   `json:"id"` // имя метрики
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type GetMetricRequest struct {
	ID    string `json:"id" binding:"required"` // имя метрики
	MType string `json:"MType" binding:"required"`
}

type MetricResponse struct {
	ID    string `json:"id"` // имя метрики
	MType string `json:"MType"`
	Delta int64  `json:"delta,omitempty"`
	Value string `json:"value,omitempty"`
}

func formatGauge(value float64) string {
	formatted := fmt.Sprintf("%.3f", value)
	parts := strings.Split(formatted, ".")
	if len(parts) == 2 {
		integerPart := parts[0]
		decimalPart := strings.TrimRight(parts[1], "0")
		if decimalPart == "" {
			return integerPart
		}
		return integerPart + "." + decimalPart
	}
	return formatted
}

func (s *Service) createGauge(name string, value float64) storage.GaugeMetric {
	metric := storage.GaugeMetric{
		ID:         int(time.Now().UnixNano()),
		MetricName: name,
		Timestamp:  time.Now().Unix(),
		Value:      value,
	}
	s.storage.AddGauge(metric)
	return metric
}

func (s *Service) createCounter(name string, delta int64) storage.CounterMetric {
	last, exists := s.storage.GetLastCounter(name)
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
	s.storage.AddCounter(metric)
	return metric
}

func (s *Service) UpdateMetric(req UpdateMetricRequest) (*MetricResponse, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("missing metrics name")
	}

	switch req.MType {
	case "gauge":
		return s.UpdateGauge(req.ID, *req.Value)
	case "counter":
		return s.UpdateCounter(req.ID, *req.Delta)
	default:
		return nil, fmt.Errorf("invalid metric type: %s", req.MType)
	}
}

func (s *Service) UpdateGauge(name string, value float64) (*MetricResponse, error) {
	metric := s.createGauge(name, value)
	valueFloat := fmt.Sprintf("%f", metric.Value)

	return &MetricResponse{
		ID:    fmt.Sprintf("%v", metric.ID),
		MType: "gauge",
		Value: valueFloat,
	}, nil
}

func (s *Service) UpdateCounter(name string, delta int64) (*MetricResponse, error) {
	s.createCounter(name, delta)
	newMetric, exist := s.storage.GetLastCounter(name)
	if !exist {
		return nil, fmt.Errorf("missing name")
	}

	newValue := newMetric.Value
	return &MetricResponse{
		ID:    fmt.Sprintf("%v", name),
		MType: "counter",
		Delta: newValue,
	}, nil
}

func (s *Service) GetMetric(req GetMetricRequest) (*MetricResponse, error) {
	switch req.MType {
	case "gauge":
		value, ok := s.getGaugeValue(req.ID)
		if !ok {
			return &MetricResponse{
				ID:    fmt.Sprintf("%v", req.ID),
				MType: "gauge",
				Value: "",
			}, nil
		}
		return &MetricResponse{
			ID:    fmt.Sprintf("%v", req.ID),
			MType: "gauge",
			Value: value,
		}, nil
	case "counter":
		value, ok := s.getCounterValue(req.ID)
		if !ok {
			return &MetricResponse{
				ID:    fmt.Sprintf("%v", req.ID),
				MType: "gauge",
				Value: "",
			}, nil
		}
		return &MetricResponse{
			ID:    fmt.Sprintf("%v", req.ID),
			MType: "counter",
			Value: value,
		}, nil
	}

	return nil, fmt.Errorf("neither gauge nor couter")
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

func (s *Service) getGaugeValue(name string) (string, bool) {
	if m, ok := s.storage.GetLastGauge(name); ok {
		return formatGauge(m.Value), true
	}
	return "", false
}

func (s *Service) getCounterValue(name string) (string, bool) {
	if m, ok := s.storage.GetLastCounter(name); ok {
		return strconv.FormatInt(m.Value, 10), true
	}
	return "", false
}

func (s *Service) ListAll() map[string][]string {
	gauges := make([]string, 0)
	for _, m := range s.storage.GetAllGauge() {
		formatted := formatGauge(m.Value)
		gauges = append(gauges, fmt.Sprintf("%d: %s - %s", m.ID, m.MetricName, formatted))
	}

	counters := make([]string, 0)
	for _, m := range s.storage.GetAllCounter() {
		counters = append(counters, fmt.Sprintf("%d: %s - %d", m.ID, m.MetricName, m.Value))
	}

	return map[string][]string{
		"gauge":   gauges,
		"counter": counters,
	}
}
