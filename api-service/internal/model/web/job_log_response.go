package web

import "time"

type JobLogResponse struct {
	Id         int64      `json:"id"`
	Status     string     `json:"status"`
	Message    string     `json:"message"`
	ExecutedAt *time.Time `json:"executed_at"`
}
