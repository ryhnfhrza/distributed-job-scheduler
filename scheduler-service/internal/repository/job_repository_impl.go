package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ryhnfhrza/distributed-job-scheduler/scheduler-service/internal/model/domain"
)

type JobRepositoryImpl struct {
	Db *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &JobRepositoryImpl{
		Db: db,
	}
}

func (repository *JobRepositoryImpl) FetchAndMarkQueued(ctx context.Context, limit int) ([]domain.Job, error) {
	tx, err := repository.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id FROM jobs
		WHERE run_at <= UTC_TIMESTAMP() AND status = 'pending'
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

	query := "UPDATE jobs SET status = 'queued',updated_at = UTC_TIMESTAMP() WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"

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
