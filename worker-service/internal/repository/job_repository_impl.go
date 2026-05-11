package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ryhnfhrza/distributed-job-scheduler/worker-service/internal/model/domain"
)

type jobRepositoryImpl struct {
	Db *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepositoryImpl{
		Db: db,
	}
}

func (repository *jobRepositoryImpl) FindById(ctx context.Context, id int64) (domain.Job, error) {
	query := "select id,name,type,payload,run_at,status,retry_count,max_retry,last_error,created_at,updated_at from jobs where id = ? "
	row := repository.Db.QueryRowContext(ctx, query, id)

	job := domain.Job{}
	err := row.Scan(
		&job.Id,
		&job.Name,
		&job.Type,
		&job.Payload,
		&job.RunAt,
		&job.Status,
		&job.RetryCount,
		&job.MaxRetry,
		&job.LastError,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return job, fmt.Errorf("job not found (id=%d): %w", id, err)
		}
		return job, fmt.Errorf("failed to scan job (id=%d): %w", id, err)
	}
	return job, nil
}

func (repository *jobRepositoryImpl) UpdateExecutionResult(ctx context.Context, job domain.Job) error {
	query := " UPDATE jobs SET status = ?, retry_count = ?, last_error = ?, run_at=?,updated_at = UTC_TIMESTAMP() WHERE id = ? AND status = 'queued' "

	result, err := repository.Db.ExecContext(ctx, query, job.Status, job.RetryCount, job.LastError, job.RunAt, job.Id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("job not found or not updated")
	}

	return nil
}

func (repository *jobRepositoryImpl) InsertLog(ctx context.Context, log domain.JobLog) error {
	query := "insert into job_logs(job_id,status,message) values(?,?,?)"
	_, err := repository.Db.ExecContext(ctx, query, log.JobId, log.Status, log.Message)
	if err != nil {
		return err
	}
	return nil
}
