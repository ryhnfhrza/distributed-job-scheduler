package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/config"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/helper"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/model/domain"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/model/web"
	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/repository"
	emailtemplate "github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/template"
)

type workerServiceImpl struct {
	JobRepository   repository.JobRepository
	QueueRepository repository.QueueRepository
	SMTPConfig      config.SMTPConfig
}

func NewWorkerService(jobRepository repository.JobRepository, queueRepository repository.QueueRepository, smtpConfig config.SMTPConfig) WorkerService {
	return &workerServiceImpl{
		JobRepository:   jobRepository,
		QueueRepository: queueRepository,
		SMTPConfig:      smtpConfig,
	}
}

func (service *workerServiceImpl) Process(ctx context.Context) error {
	data, err := service.QueueRepository.Consume(ctx)
	if err != nil {
		return err
	}

	var job domain.Job
	if err := json.Unmarshal(data, &job); err != nil {
		return err
	}

	err = service.executeJob(ctx, job)

	if err != nil {
		job.RetryCount++
		job.LastError = sql.NullString{String: err.Error(), Valid: true}

		if job.RetryCount >= job.MaxRetry {
			job.Status = "failed"
		} else {
			job.Status = "pending"
			job.Status = "pending"

			nextRun := time.Now().Add(
				time.Minute * time.Duration(job.RetryCount),
			)

			job.RunAt = &nextRun
		}
	} else {
		job.Status = "completed"
		job.LastError = sql.NullString{String: "", Valid: false}
	}

	if err := service.JobRepository.UpdateExecutionResult(ctx, job); err != nil {
		return err
	}

	message := "success"
	if job.LastError.Valid {
		message = job.LastError.String
	}

	log := domain.JobLog{
		JobId:   int(job.Id),
		Status:  job.Status,
		Message: message,
	}

	if err := service.JobRepository.InsertLog(ctx, log); err != nil {
		return err
	}

	return nil
}

func (service *workerServiceImpl) executeJob(ctx context.Context, job domain.Job) error {
	switch strings.ToUpper(job.Type) {
	case "EMAIL":
		return service.executeEmail(ctx, job)
	case "WEBHOOK":
		return service.executeWebhook(ctx, job)
	default:
		return fmt.Errorf("unsupported job type: %s", job.Type)
	}
}

func (service *workerServiceImpl) executeEmail(ctx context.Context, job domain.Job) error {
	var p web.EmailPayload
	if err := json.Unmarshal(job.Payload, &p); err != nil {
		return err
	}

	htmlBody, err := emailtemplate.RenderEmail(emailtemplate.EmailData{
		Title:   p.Subject,
		Content: p.Body,
	})
	if err != nil {
		return err
	}

	err = helper.SendEmail(service.SMTPConfig, p.To, p.Subject, htmlBody)

	if err != nil {
		return fmt.Errorf("failed to process email job: %w", err)
	}

	return nil
}

func (service *workerServiceImpl) executeWebhook(ctx context.Context, job domain.Job) error {
	var payload web.WebhookPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return err
	}

	bodyBytes, err := json.Marshal(payload.Body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.URL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	for k, v := range payload.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook failed with status: %d", resp.StatusCode)
	}

	return nil
}
