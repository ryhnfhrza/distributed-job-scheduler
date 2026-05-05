package service

import (
	"context"

	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/repository"
)

type SchedulerServiceImpl struct {
	JobRepository   repository.JobRepository
	QueueRepository repository.QueueRepository
}

func NewSchedulerService(jobRepository repository.JobRepository, queueRepository repository.QueueRepository) SchedulerService {
	return &SchedulerServiceImpl{
		JobRepository:   jobRepository,
		QueueRepository: queueRepository,
	}
}

func (service *SchedulerServiceImpl) Process(ctx context.Context) error {
	jobs, err := service.JobRepository.FetchAndMarkQueued(ctx, 10)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		if err := service.QueueRepository.Push(ctx, job); err != nil {
			return err
		}
	}

	return nil
}
