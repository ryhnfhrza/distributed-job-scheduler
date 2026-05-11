package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/exception"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/helper"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/domain"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/web"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/repository"
)

type jobServiceImpl struct {
	JobRepository repository.JobRepository
	Validate      *validator.Validate
}

func NewJobService(jobRepository repository.JobRepository, validate *validator.Validate) JobService {
	return &jobServiceImpl{
		JobRepository: jobRepository,
		Validate:      validate,
	}
}

func (service *jobServiceImpl) Create(ctx context.Context, createRequest web.JobCreateRequest) web.JobResponse {
	createRequest.Type = strings.ToLower(createRequest.Type)

	err := service.Validate.Struct(createRequest)
	helper.PanicIfError(err)

	if createRequest.RunAt.Before(time.Now()) {
		panic(exception.NewBadRequest("run_at cannot be in the past"))
	}

	payloadBytes, err := json.Marshal(createRequest.Payload)
	helper.PanicIfError(err)

	switch strings.ToUpper(createRequest.Type) {
	case "EMAIL":
		var p web.EmailPayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			panic("invalid email payload format")
		}
		helper.PanicIfError(service.Validate.Struct(p))

	case "WEBHOOK":
		var p web.WebhookPayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			panic("invalid webhook payload format")
		}
		helper.PanicIfError(service.Validate.Struct(p))

	default:
		panic("unsupported job type")
	}

	userRunAt := createRequest.RunAt.UTC()
	if userRunAt.Before(time.Now().UTC()) {
		panic(exception.NewBadRequest("run_at cannot be in the past"))
	}

	job := domain.Job{
		Name:    createRequest.Name,
		Type:    createRequest.Type,
		RunAt:   &userRunAt,
		Payload: payloadBytes,
	}

	job, err = service.JobRepository.Create(ctx, job)
	helper.PanicIfError(err)

	return helper.ToJobResponse(job)
}

func (service *jobServiceImpl) Update(ctx context.Context, updateRequest web.JobUpdateRequest) web.JobResponse {
	err := service.Validate.Struct(updateRequest)
	helper.PanicIfError(err)

	job, err := service.JobRepository.FindById(ctx, updateRequest.Id)
	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("job with id:%d not found", updateRequest.Id)))
	}
	if job.Status != "pending" {
		panic(exception.NewBadRequest("only pending job can be updated"))
	}

	if updateRequest.Name != "" {
		job.Name = updateRequest.Name
	}
	fmt.Println(job.RunAt)

	if updateRequest.RunAt != nil {
		if !updateRequest.RunAt.IsZero() {
			utcRunAt := updateRequest.RunAt.UTC()
			job.RunAt = &utcRunAt
		}
	}

	job, err = service.JobRepository.Update(ctx, job)
	helper.PanicIfError(err)

	return helper.ToJobResponse(job)
}

func (service *jobServiceImpl) Delete(ctx context.Context, id int64) {

	job, err := service.JobRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("job with id:%d not found", id)))
	}

	if job.Status != "pending" {
		panic(exception.NewBadRequest("only pending job can be updated"))
	}

	err = service.JobRepository.Delete(ctx, id)
	helper.PanicIfError(err)
}

func (service *jobServiceImpl) FindById(ctx context.Context, id int64) web.JobResponse {
	job, err := service.JobRepository.FindById(ctx, id)
	if err != nil {
		panic(exception.NewNotFoundError(fmt.Sprintf("job with id:%d not found", id)))
	}

	return helper.ToJobResponse(job)
}

func (service *jobServiceImpl) FindAll(ctx context.Context, filter domain.Filter) []web.JobResponse {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	jobs, err := service.JobRepository.FindAll(ctx, filter)
	helper.PanicIfError(err)

	return helper.ToJobResponses(jobs)
}

func (service *jobServiceImpl) FindJobLog(ctx context.Context, idJob int64) []web.JobLogResponse {
	jobLogs, err := service.JobRepository.GetLogs(ctx, idJob)
	helper.PanicIfError(err)

	return helper.ToJobLogResponses(jobLogs)
}
