package service

import "context"

type WorkerService interface {
	Process(ctx context.Context) error
}
