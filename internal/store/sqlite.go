package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/o0ga-bo0ga/vigil/internal/models"
	_ "github.com/mattn/go-sqlite3"
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

func (s *SQLiteStore) ListJobs(tenant string) ([]*models.Job, error) {
	var jobs []*models.Job
	var rows *sql.Rows
	var err error
	query := `select id, name, status, error, duration, tenant, created_at, updated_at from jobs`
	if tenant != "" {
		query += ` where tenant = ?`
		rows, err = s.db.Query(query, tenant)
	} else {
		rows, err = s.db.Query(query)
	}
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
		return errors.New("Invalid job ID")
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