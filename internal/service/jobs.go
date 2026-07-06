package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/o0ga-bo0ga/vigil/internal/models"
	"github.com/o0ga-bo0ga/vigil/internal/store"
)

type CreateJobRequest struct {
	Name      string
	Status    models.Status
	Error     string
	Duration  int64
	Tenant    string
}

type JobService struct {
	store store.Store
}

type UpdateJobRequest struct {
	ID       string
	Status   models.Status
	Error    string
	Duration int64
}

func NewJobService(store store.Store) *JobService {
	return &JobService{store: store}
}

func (s *JobService) CreateJob(req CreateJobRequest) (*models.Job, error) {
	id := uuid.New().String()
	
	status := req.Status
	if status == "" {
		status = models.StatusStarted
	}

	job := models.Job{
		ID: id,
		Name: req.Name,
		Status: status,
		Tenant: req.Tenant,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.store.CreateJob(&job)
	return &job, err
}

func (s *JobService) GetJob(id string) (*models.Job, error) {
	return s.store.GetJob(id)
}

func (s *JobService) ListJobs(tenant string) ([]*models.Job, error) {
	return s.store.ListJobs(tenant)
}

func (s *JobService) UpdateJob(req UpdateJobRequest) error {
	job, err := s.store.GetJob(req.ID)

	if err != nil {
		return err
	}

	job.Status = req.Status
	job.Error = req.Error
	job.Duration = req.Duration
	job.UpdatedAt = time.Now()
	return s.store.UpdateJob(job)
}