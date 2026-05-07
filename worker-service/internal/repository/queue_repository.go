package repository

import "context"

type QueueRepository interface {
	Consume(ctx context.Context) ([]byte, error)
}
