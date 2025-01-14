package metrics

import (
	"metrics-persistance-server/internal/message"
	"metrics-persistance-server/internal/metrics/repo"
	"time"
)

type Service struct {
	repo repo.IRepository
}

type IService interface {
	RecordMetric(message *message.Message) error
	GetMetricsByRoute(route string, startTime, endTime time.Time) ([]repo.ApiMetric, error)
}

func NewApiMetricService(repo repo.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordMetric(message *message.Message) error {
	metric := &repo.ApiMetric{
		Time:         time.Unix(int64(message.Timestamp), 0),
		Route:        message.Payload.Route,
		Method:       message.Payload.Method,
		StatusCode:   message.Payload.StatusCode,
		ResponseTime: message.Payload.ResponseTime,
	}
	return s.repo.Insert(metric)
}

func (s *Service) GetMetricsByRoute(route string, startTime, endTime time.Time) ([]repo.ApiMetric, error) {
	return s.repo.GetByRouteAndTime(route, startTime, endTime)
}
