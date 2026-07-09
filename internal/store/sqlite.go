package store

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/o0ga-bo0ga/vigil/internal/models"
)

type SQLiteStore struct {
	db *sql.DB
}

var _ Store = (*SQLiteStore)(nil)

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	migration := `Create Table if not exists jobs(
					id text primary key,
					name text not null,
					status text not null,
					error text default '',
					duration integer default 0,
					tenant text not null,
					created_at datetime not null,
					updated_at datetime not null)`

	_, err = db.Exec(migration)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) CreateJob(job *models.Job) error {
	if job == nil || job.ID == ""{
		return errors.New("invalid job ID")
	}
	query := `insert into jobs (id,
								name,
								status,
								error,
								duration,
								tenant, 
								created_at, 
								updated_at)
							values (?,?,?,?,?,?,?,?)`
	_, err := s.db.Exec(query,
				job.ID,
				job.Name,
				job.Status,
				job.Error,
				job.Duration,
				job.Tenant,
				job.CreatedAt,
				job.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) GetJob(id string) (*models.Job, error) {
	if id == "" {
		return nil, ErrNotFound
	}
	var job models.Job
	query := `select id, name, status, error, duration, tenant, created_at, updated_at from jobs where id = ?`
	if err := s.db.QueryRow(query, id).Scan(&job.ID,
											&job.Name,
											&job.Status,
											&job.Error,
											&job.Duration,
											&job.Tenant,
											&job.CreatedAt,
											&job.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &job, nil
}

func (s *SQLiteStore) ListJobs(filter ListJobsFilter) ([]*models.Job, error) {
	var jobs []*models.Job
	var rows *sql.Rows
	var err error
	query := `select id, name, status, error, duration, tenant, created_at, updated_at from jobs`
	
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	} 
	if filter.Tenant != "" {
		conditions = append(conditions, "tenant = ?")
		args = append(args, filter.Tenant)
	}
	if len(conditions) > 0 {
		query += " where " + strings.Join(conditions, " and ")
	}

	query += ` order by created_at desc`

	if filter.Limit > 0 {
		query += ` limit ?`
		args = append(args, filter.Limit)
	}
	rows, err = s.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var job models.Job
		err = rows.Scan(&job.ID,
						&job.Name,
						&job.Status,
						&job.Error,
						&job.Duration,
						&job.Tenant,
						&job.CreatedAt,
						&job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *SQLiteStore) UpdateJob(job *models.Job) error {
	if job == nil || job.ID == "" {
		return errors.New("invalid job ID")
	}
	query := `update jobs set status = ?, error = ?, duration = ?, updated_at = ? where id = ?`
	result, err := s.db.Exec(query, job.Status, job.Error, job.Duration, time.Now(), job.ID)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) GetStats(tenant string) (*Stats, error) {
	var stats Stats
	var args []interface{}
	query := `select 
				count(*) as Total, 
				sum(status = 'succeeded') as Succeeded,
				sum(status = 'failed') as Failed,
				sum(status = 'retried') as Retried,
				sum(status = 'started') as Started,
				coalesce(avg(duration), 0) as AvgDuration
				from jobs`
	if tenant != "" {
		query += ` where tenant = ?`
		args = append(args, tenant)
	}
	err := s.db.QueryRow(query, args...).Scan(&stats.Total,
									          &stats.Succeeded,
										      &stats.Failed,
										      &stats.Retried,
										      &stats.Started,
										      &stats.AvgDuration)
	if err != nil {
		return nil, err
	}

	return &stats, nil

}