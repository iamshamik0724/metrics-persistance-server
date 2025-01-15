package metrics

import (
	"fmt"
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
	FetchLast10MinutesMetrics() (*Metrics, error)
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

func (s *Service) FetchLast10MinutesMetrics() (*Metrics, error) {
	endTime := time.Now()
	startTime := endTime.Add(-10 * time.Minute)

	apiMetrics, err := s.repo.GetByTimeBetween(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error while fetching metrics: %w", err)
	}

	metricsResponse := &Metrics{
		MetricsData: []RouteMetric{},
	}

	routeMetricsMap := make(map[string]RouteMetric)

	for _, apiMetric := range apiMetrics {
		routeKey := fmt.Sprintf("%s|%s", apiMetric.Route, apiMetric.Method)
		if _, exists := routeMetricsMap[routeKey]; !exists {
			routeMetricsMap[routeKey] = RouteMetric{
				Route:      apiMetric.Route,
				Method:     apiMetric.Method,
				Timestamps: []time.Time{},
				Responses:  []ResponseMetric{},
			}
		}

		routeMetric := routeMetricsMap[routeKey]
		routeMetric.Timestamps = append(routeMetric.Timestamps, apiMetric.Time)
		routeMetric.Responses = append(routeMetric.Responses, ResponseMetric{
			Time:   apiMetric.ResponseTime,
			Status: apiMetric.StatusCode,
		})
		routeMetricsMap[routeKey] = routeMetric
	}

	for _, routeMetric := range routeMetricsMap {
		metricsResponse.MetricsData = append(metricsResponse.MetricsData, routeMetric)
	}

	return metricsResponse, nil
}
