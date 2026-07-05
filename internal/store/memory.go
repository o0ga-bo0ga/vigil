package store

import (
	"errors"
	"sync"

	"github.com/o0ga-bo0ga/vigil/internal/models"
)

var _ Store = (*MemoryStore)(nil)

type MemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]*models.Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]*models.Job),
	}
}

func (s *MemoryStore) CreateJob(job *models.Job) error {
	if job == nil || job.ID == "" {
		return errors.New("Invalid Job ID")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; exists {
		return errors.New("job already exists")
	}

	jobCopy := *job
	s.jobs[job.ID] = &jobCopy
	return nil
}

func (s *MemoryStore) GetJob(id string) (*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[id]
	if !exists {
		return nil, errors.New("job not found")
	}

	jobCopy := *job
	return &jobCopy, nil
}

func (s *MemoryStore) ListJobs(tenant string) ([]*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tenantJobs []*models.Job
	for _, job := range s.jobs {
		if tenant == "" || job.Tenant == tenant {
			jobCopy := *job
			tenantJobs = append(tenantJobs, &jobCopy)
		}
	}

	return tenantJobs, nil
}

func (s *MemoryStore) UpdateJob(job *models.Job) error {
	if job == nil || job.ID == "" {
		return errors.New("Invalid Job ID")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; !exists {
		return errors.New("job not found")
	}

	jobCopy := *job
	s.jobs[job.ID] = &jobCopy
	return nil
}