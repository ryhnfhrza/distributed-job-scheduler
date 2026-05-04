package web

import "time"

type JobResponse struct {
	Id         int64       `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Payload    interface{} `json:"payload"`
	RunAt      *time.Time  `json:"run_at"`
	Status     string      `json:"status,omitempty"`
	RetryCount int         `json:"retry_count,omitempty"`
	MaxRetry   int         `json:"max_retry,omitempty"`
	LastError  string      `json:"last_error,omitempty"`
	CreatedAt  *time.Time  `json:"created_at,omitempty"`
	UpdateAt   *time.Time  `json:"updated_at,omitempty"`
}
