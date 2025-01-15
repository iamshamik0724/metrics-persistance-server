package metrics

import "time"

type Metrics struct {
	MetricsData []RouteMetric `json:"metrics"`
}

type RouteMetric struct {
	Route      string           `json:"route"`
	Method     string           `json:"method"`
	Timestamps []time.Time      `json:"timestamps"`
	Responses  []ResponseMetric `json:"responses"`
}

type ResponseMetric struct {
	Time   float64 `json:"time"`
	Status int     `json:"status"`
}
