package domain

import (
	"database/sql"
	"time"
)

type Job struct {
	Id         int64
	Name       string
	Type       string
	Payload    []byte
	RunAt      *time.Time
	Status     string
	RetryCount int
	MaxRetry   int
	LastError  sql.NullString
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
