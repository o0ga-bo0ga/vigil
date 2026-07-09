package store

import (
	"errors"

	"github.com/o0ga-bo0ga/vigil/internal/models"
)

type Store interface {
	CreateJob(job *models.Job) error
	GetJob(id string) (*models.Job, error)
	ListJobs(filter ListJobsFilter) ([]*models.Job, error)
	UpdateJob(job *models.Job) error
	GetStats(tenant string) (*Stats, error)
}

type ListJobsFilter struct {
	Tenant string
	Status string
	Limit  int
}

type Stats struct {
	Total       int
	Succeeded   int
	Failed      int
	Retried     int
	Started     int
	AvgDuration float64
}

var (
	ErrNotFound = errors.New("job not found")
)