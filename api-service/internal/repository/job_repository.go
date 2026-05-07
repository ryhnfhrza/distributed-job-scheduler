package repository

import (
	"context"

	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job domain.Job) (domain.Job, error)
	Update(ctx context.Context, job domain.Job) (domain.Job, error)
	Delete(ctx context.Context, id int64) error
	FindById(ctx context.Context, id int64) (domain.Job, error)
	FindAll(ctx context.Context, filter domain.Filter) ([]domain.Job, error)
	GetLogs(ctx context.Context, idJob int64) ([]domain.JobLog, error)
}
