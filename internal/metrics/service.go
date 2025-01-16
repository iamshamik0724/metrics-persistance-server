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
		Time:         message.Timestamp,
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

	responseTimestamps := generateTimeArray(startTime, endTime)
	apiMetrics, err := s.repo.GetByTimeBetween(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error while fetching metrics: %w", err)
	}

	metricsResponse := &Metrics{
		Timestamps:  responseTimestamps,
		MetricsData: []RouteMetric{},
	}

	routeTimestampMetricsMap := make(map[string]map[time.Time]ResponseMetric)

	for _, apiMetric := range apiMetrics {
		routeKey := fmt.Sprintf("%s - %s", apiMetric.Method, apiMetric.Route)

		if _, exists := routeTimestampMetricsMap[routeKey]; !exists {
			routeTimestampMetricsMap[routeKey] = make(map[time.Time]ResponseMetric)
		}
		routeTimestampMetricsMap[routeKey][apiMetric.Time] = ResponseMetric{
			Time:   apiMetric.ResponseTime,
			Status: apiMetric.StatusCode,
		}
	}

	for routeKey, timeMetricMap := range routeTimestampMetricsMap {
		responseTimes := make([]float64, len(responseTimestamps))
		responseStatus := make([]int, len(responseTimestamps))

		for i := 0; i < len(metricsResponse.Timestamps); i++ {
			if res, exists := timeMetricMap[metricsResponse.Timestamps[i]]; exists {
				responseTimes[i] = res.Time
				responseStatus[i] = res.Status
			} else {
				responseTimes[i] = -1
				responseStatus[i] = -1
			}
		}
		metricsResponse.MetricsData = append(metricsResponse.MetricsData, RouteMetric{
			RouteKey:     routeKey,
			ResponseTime: responseTimes,
			Status:       responseStatus,
		})
	}

	return metricsResponse, nil
}

func generateTimeArray(startTime, endTime time.Time) []time.Time {
	var timeArray []time.Time
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)
	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(1 * time.Second) {
		timeArray = append(timeArray, t)
	}
	return timeArray
}
