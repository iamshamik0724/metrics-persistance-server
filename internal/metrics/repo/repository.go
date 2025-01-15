package repo

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type IRepository interface {
	Insert(metric *ApiMetric) error
	GetByRouteAndTime(route string, startTime, endTime time.Time) ([]ApiMetric, error)
	GetByTimeBetween(startTime, endTime time.Time) ([]ApiMetric, error)
}

func NewApiMetricRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(metric *ApiMetric) error {
	return r.db.Create(metric).Error
}

func (r *Repository) GetByRouteAndTime(route string, startTime, endTime time.Time) ([]ApiMetric, error) {
	var metrics []ApiMetric
	err := r.db.Where("route = ? AND time BETWEEN ? AND ?", route, startTime, endTime).Find(&metrics).Error
	return metrics, err
}

func (r *Repository) GetByTimeBetween(startTime, endTime time.Time) ([]ApiMetric, error) {
	var metrics []ApiMetric
	err := r.db.Where("time BETWEEN ? AND ?", startTime, endTime).Find(&metrics).Error
	return metrics, err
}
