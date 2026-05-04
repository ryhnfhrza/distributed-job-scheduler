package service

import (
	"context"

	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/domain"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/web"
)

type JobService interface {
	Create(ctx context.Context, createRequest web.JobCreateRequest) web.JobResponse
	Update(ctx context.Context, updateRequest web.JobUpdateRequest) web.JobResponse
	Delete(ctx context.Context, id int64)
	FindById(ctx context.Context, id int64) web.JobResponse
	FindAll(ctx context.Context, filter domain.Filter) []web.JobResponse
	FindJobLog(ctx context.Context, idJob int64) []web.JobLogResponse
}
