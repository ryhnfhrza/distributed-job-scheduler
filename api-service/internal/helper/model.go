package helper

import (
	"encoding/json"

	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/domain"
	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/web"
)

func ToJobResponse(job domain.Job) web.JobResponse {
	var decodedPayload interface{}

	if len(job.Payload) > 0 {
		err := json.Unmarshal(job.Payload, &decodedPayload)
		if err != nil {
			decodedPayload = string(job.Payload)
		}
	}
	return web.JobResponse{
		Id:         job.Id,
		Name:       job.Name,
		Type:       job.Type,
		Payload:    decodedPayload,
		RunAt:      job.RunAt,
		Status:     job.Status,
		RetryCount: job.RetryCount,
		MaxRetry:   job.MaxRetry,
		LastError:  job.LastError,
		CreatedAt:  job.CreatedAt,
		UpdateAt:   job.UpdatedAt,
	}
}

func ToJobResponses(jobs []domain.Job) []web.JobResponse {
	var jobResponses []web.JobResponse
	for _, job := range jobs {
		jobResponses = append(jobResponses, ToJobResponse(job))
	}
	return jobResponses
}

func ToJobLogResponse(jobLog domain.JobLog) web.JobLogResponse {
	return web.JobLogResponse{
		Id:         jobLog.Id,
		Status:     jobLog.Status,
		Message:    jobLog.Message,
		ExecutedAt: jobLog.ExecutedAt,
	}
}

func ToJobLogResponses(jobLogs []domain.JobLog) []web.JobLogResponse {
	var jobLogResponses []web.JobLogResponse
	for _, jobLog := range jobLogs {
		jobLogResponses = append(jobLogResponses, ToJobLogResponse(jobLog))
	}
	return jobLogResponses
}
