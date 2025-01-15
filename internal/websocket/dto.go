package websocket

import "time"

type MetricsBroadcastMessage struct {
	Route        string    `json:"route"`
	Method       string    `json:"method"`
	Timestamp    time.Time `json:"timestamp"`
	ResponseTime float64   `json:"time"`
	Status       int       `json:"status"`
}
