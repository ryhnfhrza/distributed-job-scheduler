package repository

import (
	"context"

	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/model/domain"
)

type JobRepository interface {
	FetchAndMarkQueued(ctx context.Context, limit int) ([]domain.Job, error)
}
