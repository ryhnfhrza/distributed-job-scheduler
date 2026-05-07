package repository

import (
	"context"

	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/model/domain"
)

type JobRepository interface {
	FindById(ctx context.Context, id int64) (domain.Job, error)
	UpdateExecutionResult(ctx context.Context, job domain.Job) error
	InsertLog(ctx context.Context, log domain.JobLog) error
}
