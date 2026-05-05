package service

import "context"

type SchedulerService interface {
	Process(ctx context.Context) error
}
