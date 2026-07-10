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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Error     string    `json:"error"`
	Duration  int64     `json:"duration"`
	Tenant    string    `json:"tenant"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}