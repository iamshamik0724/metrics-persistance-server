package repo

import "time"

type ApiMetric struct {
	Time         time.Time `gorm:"primaryKey;column:time"`
	Route        string    `gorm:"primaryKey;column:route"`
	Method       string    `gorm:"primaryKey;column:method"`
	StatusCode   int       `gorm:"column:status_code"`
	ResponseTime float64   `gorm:"column:response_time"`
}

func (ApiMetric) TableName() string {
	return "api_metrics"
}
