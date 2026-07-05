package store

import (
	"errors"

	"github.com/o0ga-bo0ga/vigil/internal/models"
)

type Store interface {
	CreateJob(job *models.Job) error
	GetJob(id string) (*models.Job, error)
	ListJobs(tenant string) ([]*models.Job, error)
	UpdateJob(job *models.Job) error
}

var (
	ErrNotFound = errors.New("job not found")
)