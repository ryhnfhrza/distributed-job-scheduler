package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type queueRepositoryImpl struct {
	RedisClient *redis.Client
}

func NewQueueRepository(redisClient *redis.Client) QueueRepository {
	return &queueRepositoryImpl{
		RedisClient: redisClient,
	}
}

func (repository *queueRepositoryImpl) Consume(ctx context.Context) ([]byte, error) {
	result, err := repository.RedisClient.BLPop(ctx, 0, "job_queue").Result()
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid redis response")
	}

	return []byte(result[1]), nil
}
