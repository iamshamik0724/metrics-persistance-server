package metrics

import "time"

type Metrics struct {
	Timestamps  []time.Time   `json:"timestamps"`
	MetricsData []RouteMetric `json:"metrics"`
}

type RouteMetric struct {
	RouteKey     string    `json:"route"`
	ResponseTime []float64 `json:"responseTime"`
	Status       []int     `json:"responseStatus"`
}

type ResponseMetric struct {
	Time   float64 `json:"time"`
	Status int     `json:"status"`
}
