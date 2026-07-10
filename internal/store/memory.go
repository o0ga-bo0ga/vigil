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
		return nil, ErrNotFound
	}

	jobCopy := *job
	return &jobCopy, nil
}

func (s *MemoryStore) ListJobs(filter ListJobsFilter) ([]*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filteredJobs := make([]*models.Job, 0)
	for _, job := range s.jobs {
		if (filter.Tenant == "" || job.Tenant == filter.Tenant) &&
		   (filter.Status == "" || job.Status == models.Status(filter.Status)) {
			jobCopy := *job
			filteredJobs = append(filteredJobs, &jobCopy)
		}
	}
	if filter.Limit > 0 && len(filteredJobs) > filter.Limit {
		filteredJobs = filteredJobs[:filter.Limit]
	}

	return filteredJobs, nil
}

func (s *MemoryStore) UpdateJob(job *models.Job) error {
	if job == nil || job.ID == "" {
		return errors.New("invalid job ID")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; !exists {
		return ErrNotFound
	}

	jobCopy := *job
	s.jobs[job.ID] = &jobCopy
	return nil
}

func (s *MemoryStore) GetStats(tenant string) (*Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var duration int64
	
	var stats Stats

	for _, job := range s.jobs {
		if tenant != "" && job.Tenant != tenant {
			continue
		}
		stats.Total++
		if job.Status == models.StatusSucceeded {
			stats.Succeeded++
		} else if job.Status == models.StatusFailed {
			stats.Failed++
		} else if job.Status == models.StatusRetried {
			stats.Retried++
		} else if job.Status == models.StatusStarted {
			stats.Started++
		}
		duration += job.Duration
	}
	if stats.Total > 0 {
		stats.AvgDuration = float64(duration)/float64(stats.Total)
	} else {
		stats.AvgDuration = 0
	}

	return &stats, nil
}