package domain

import "time"

type JobLog struct {
	Id         int64
	JobId      int
	Status     string
	Message    string
	ExecutedAt *time.Time
}
