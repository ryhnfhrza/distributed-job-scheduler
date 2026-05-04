package domain

import "time"

type Job struct {
	Id         int64
	Name       string
	Type       string
	Payload    []byte
	RunAt      *time.Time
	Status     string
	RetryCount int
	MaxRetry   int
	LastError  string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
