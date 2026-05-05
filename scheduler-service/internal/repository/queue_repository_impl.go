package repository

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/model/domain"
)

type QueueRepositoryImpl struct {
	Client *redis.Client
}

func NewQueueRepository(client *redis.Client) QueueRepository {
	return &QueueRepositoryImpl{
		Client: client,
	}
}

func (r *QueueRepositoryImpl) Push(ctx context.Context, job domain.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.Client.RPush(ctx, "job_queue", data).Err()
}
