package models

import "time"

type Status string

const (
	StatusStarted   Status = "started"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusRetried   Status = "retried"
)

type Job struct {
	ID        string
	Name      string
	Status    Status
	Error     string
	Duration  int64
	Tenant    string
	CreatedAt time.Time
	UpdatedAt time.Time
}