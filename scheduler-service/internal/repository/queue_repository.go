package repository

import (
	"context"

	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/model/domain"
)

type QueueRepository interface {
	Push(ctx context.Context, job domain.Job) error
}
