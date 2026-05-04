package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ryhnfhrza/distributed-job-scheduler/api-service/internal/model/domain"
)

type JobRepositoryImpl struct {
	Db *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &JobRepositoryImpl{
		Db: db,
	}
}

func (repository *JobRepositoryImpl) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	query := "insert into jobs(name,type,payload,run_at,last_error) values (?,?,?,?,?)"
	result, err := repository.Db.ExecContext(ctx, query, job.Name, job.Type, job.Payload, job.RunAt, job.LastError)
	if err != nil {
		return job, fmt.Errorf("failed to insert job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return job, fmt.Errorf("failed to get last insert id : %w", err)
	}

	job.Id = id

	return job, nil
}

func (repository *JobRepositoryImpl) Update(ctx context.Context, job domain.Job) (domain.Job, error) {
	query := "update jobs set name = ?, run_at = ?, updated_at = NOW() where id = ? and status = 'pending'"
	result, err := repository.Db.ExecContext(ctx, query, job.Name, job.RunAt, job.Id)
	if err != nil {
		return job, fmt.Errorf("failed to update jobs (id=%d): %w", job.Id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return job, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return job, fmt.Errorf("job cannot be updated (not found or not pending)")
	}

	return job, nil
}

func (repository *JobRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := "delete from jobs where id = ? and status = 'pending'"

	result, err := repository.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job (id=%d): %w", id, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("job not found or cannot be deleted")
	}

	return nil
}

func (repository *JobRepositoryImpl) FindById(ctx context.Context, id int64) (domain.Job, error) {
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

func (repository *JobRepositoryImpl) FindAll(ctx context.Context, filter domain.Filter) ([]domain.Job, error) {
	var (
		sb   strings.Builder
		args []interface{}
	)

	sb.WriteString(`select id,name,type,payload,run_at,status,retry_count,max_retry,last_error,created_at,updated_at from jobs`)

	if filter.Status != "" {
		sb.WriteString(" WHERE status = ?")
		args = append(args, filter.Status)
	}

	if filter.Limit > 0 {
		sb.WriteString(" LIMIT ?")
		args = append(args, filter.Limit)
		if filter.Offset > 0 {
			sb.WriteString(" OFFSET ?")
			args = append(args, filter.Offset)
		}
	}
	rows, err := repository.Db.QueryContext(ctx, sb.String(), args...)

	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}

	defer rows.Close()

	jobs := []domain.Job{}

	for rows.Next() {
		job := domain.Job{}
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan job row : %w", err)
		}
		jobs = append(jobs, job)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error : %w", err)
	}

	return jobs, nil
}

func (repository *JobRepositoryImpl) FindPendingJobs(ctx context.Context) ([]domain.Job, error) {
	query := "SELECT id,name,type,payload,run_at,status,retry_count,max_retry,last_error,created_at,updated_at FROM jobs WHERE run_at <= NOW()  AND status = 'pending' LIMIT 10"
	rows, err := repository.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}

	defer rows.Close()

	jobs := []domain.Job{}

	for rows.Next() {
		job := domain.Job{}
		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan job row : %w", err)
		}
		jobs = append(jobs, job)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error : %w", err)
	}

	return jobs, nil
}

func (repository *JobRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := "update jobs set status = ?, updated_at = NOW() where id = ?"
	_, err := repository.Db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update job status (id=%d): %w", id, err)
	}

	return nil
}

func (repository *JobRepositoryImpl) GetLogs(ctx context.Context, idJob int64) ([]domain.JobLog, error) {
	query := "select id,status,message,executed_at from job_logs where job_id = ? ORDER BY executed_at DESC"
	rows, err := repository.Db.QueryContext(ctx, query, idJob)

	if err != nil {
		return nil, fmt.Errorf("failed to query job_logs: %w", err)
	}

	defer rows.Close()

	jobLogs := []domain.JobLog{}

	for rows.Next() {
		jobLog := domain.JobLog{}
		err := rows.Scan(
			&jobLog.Id,
			&jobLog.Status,
			&jobLog.Message,
			&jobLog.ExecutedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job_logs row : %w", err)
		}
		jobLogs = append(jobLogs, jobLog)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error : %w", err)
	}

	return jobLogs, nil
}

func (repository *JobRepositoryImpl) FetchAndMarkQueued(ctx context.Context, limit int) ([]domain.Job, error) {
	tx, err := repository.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id FROM jobs
		WHERE run_at <= NOW() AND status = 'pending'
		ORDER BY run_at ASC
		LIMIT ?
	`, limit)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return []domain.Job{}, nil
	}

	query := "UPDATE jobs SET status = 'queued',updated_at = NOW() WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"

	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	selectQuery := "SELECT id,name,type,payload,run_at,status,retry_count,max_retry,last_error,created_at,updated_at FROM jobs WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"

	rows2, err := tx.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows2.Close()

	var jobs []domain.Job
	for rows2.Next() {
		var job domain.Job
		if err := rows2.Scan(
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
		); err != nil {
			tx.Rollback()
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return jobs, nil
}
