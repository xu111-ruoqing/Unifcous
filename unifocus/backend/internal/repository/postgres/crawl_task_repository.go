package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/unifocus/backend/internal/domain"
)

// CrawlTaskRepository handles crawl task data access operations
type CrawlTaskRepository struct {
	db *DB
}

// NewCrawlTaskRepository creates a new crawl task repository
func NewCrawlTaskRepository(db *DB) *CrawlTaskRepository {
	return &CrawlTaskRepository{db: db}
}

// Create creates a new crawl task
func (r *CrawlTaskRepository) Create(ctx context.Context, task *domain.CrawlTask) error {
	query := `
		INSERT INTO crawl_tasks (target_url, site_name, selector_config, frequency, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	selectorConfigJSON, err := json.Marshal(task.SelectorConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal selector config: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		task.TargetURL,
		task.SiteName,
		selectorConfigJSON,
		task.Frequency,
		task.Status,
		task.CreatedAt,
	).Scan(&task.ID, &task.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create crawl task: %w", err)
	}

	return nil
}

// GetByID retrieves a crawl task by ID
func (r *CrawlTaskRepository) GetByID(ctx context.Context, id int64) (*domain.CrawlTask, error) {
	query := `
		SELECT id, target_url, site_name, selector_config, frequency, last_crawled_at, next_crawl_at, status, error_message, created_at
		FROM crawl_tasks
		WHERE id = $1
	`

	task := &domain.CrawlTask{}
	var selectorConfigJSON []byte
	var lastCrawledAt, nextCrawlAt sql.NullTime
	var errMsg sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.TargetURL,
		&task.SiteName,
		&selectorConfigJSON,
		&task.Frequency,
		&lastCrawledAt,
		&nextCrawlAt,
		&task.Status,
		&errMsg,
		&task.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("crawl task not found")
		}
		return nil, fmt.Errorf("failed to get crawl task: %w", err)
	}

	if len(selectorConfigJSON) > 0 {
		if err := json.Unmarshal(selectorConfigJSON, &task.SelectorConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal selector config: %w", err)
		}
	}

	if lastCrawledAt.Valid {
		t := lastCrawledAt.Time
		task.LastCrawledAt = &t
	}
	if nextCrawlAt.Valid {
		t := nextCrawlAt.Time
		task.NextCrawlAt = &t
	}
	if errMsg.Valid {
		task.ErrorMessage = errMsg.String
	} else {
		task.ErrorMessage = ""
	}

	return task, nil
}

// List retrieves a list of crawl tasks with pagination
func (r *CrawlTaskRepository) List(ctx context.Context, offset, limit int) ([]*domain.CrawlTask, error) {
	query := `
		SELECT id, target_url, site_name, selector_config, frequency, last_crawled_at, next_crawl_at, status, error_message, created_at
		FROM crawl_tasks
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list crawl tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.CrawlTask
	for rows.Next() {
		task := &domain.CrawlTask{}
		var selectorConfigJSON []byte
		var lastCrawledAt, nextCrawlAt sql.NullTime
		var errMsg sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.TargetURL,
			&task.SiteName,
			&selectorConfigJSON,
			&task.Frequency,
			&lastCrawledAt,
			&nextCrawlAt,
			&task.Status,
			&errMsg,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crawl task: %w", err)
		}

		if len(selectorConfigJSON) > 0 {
			if err := json.Unmarshal(selectorConfigJSON, &task.SelectorConfig); err != nil {
				return nil, fmt.Errorf("failed to unmarshal selector config: %w", err)
			}
		}

		if lastCrawledAt.Valid {
			t := lastCrawledAt.Time
			task.LastCrawledAt = &t
		}
		if nextCrawlAt.Valid {
			t := nextCrawlAt.Time
			task.NextCrawlAt = &t
		}
		if errMsg.Valid {
			task.ErrorMessage = errMsg.String
		} else {
			task.ErrorMessage = ""
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates a crawl task
func (r *CrawlTaskRepository) Update(ctx context.Context, task *domain.CrawlTask) error {
	query := `
		UPDATE crawl_tasks
		SET target_url = $1, site_name = $2, selector_config = $3, frequency = $4, last_crawled_at = $5, next_crawl_at = $6, status = $7, error_message = $8
		WHERE id = $9
	`

	selectorConfigJSON, err := json.Marshal(task.SelectorConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal selector config: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		task.TargetURL,
		task.SiteName,
		selectorConfigJSON,
		task.Frequency,
		task.LastCrawledAt,
		task.NextCrawlAt,
		task.Status,
		task.ErrorMessage,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update crawl task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("crawl task not found")
	}

	return nil
}

// Delete deletes a crawl task by ID
func (r *CrawlTaskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM crawl_tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete crawl task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("crawl task not found")
	}

	return nil
}

// GetDueTasks retrieves tasks that are due for crawling
func (r *CrawlTaskRepository) GetDueTasks(ctx context.Context, limit int) ([]*domain.CrawlTask, error) {
	query := `
		SELECT id, target_url, site_name, selector_config, frequency, last_crawled_at, next_crawl_at, status, error_message, created_at
		FROM crawl_tasks
		WHERE status != 'stopped' AND (next_crawl_at IS NULL OR next_crawl_at <= NOW())
		ORDER BY next_crawl_at ASC NULLS FIRST
		LIMIT $1
	`
	// Note: 'next_crawl_at ASC NULLS FIRST' ensures new tasks (NULL) or overdue tasks are picked first.
	// Postgres default for ASC is NULLS LAST, so we might need NULLS FIRST.
	// Let's assume standard behavior first or explicit.
	// Actually, NULL usually means "Run immediately" for us if we initialize it that way, or we can treat NULL as "never run yet".

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get due crawl tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.CrawlTask
	for rows.Next() {
		task := &domain.CrawlTask{}
		var selectorConfigJSON []byte
		var lastCrawledAt, nextCrawlAt sql.NullTime
		var errMsg sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.TargetURL,
			&task.SiteName,
			&selectorConfigJSON,
			&task.Frequency,
			&lastCrawledAt,
			&nextCrawlAt,
			&task.Status,
			&errMsg,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crawl task: %w", err)
		}

		if len(selectorConfigJSON) > 0 {
			if err := json.Unmarshal(selectorConfigJSON, &task.SelectorConfig); err != nil {
				return nil, fmt.Errorf("failed to unmarshal selector config: %w", err)
			}
		}

		if lastCrawledAt.Valid {
			t := lastCrawledAt.Time
			task.LastCrawledAt = &t
		}
		if nextCrawlAt.Valid {
			t := nextCrawlAt.Time
			task.NextCrawlAt = &t
		}
		if errMsg.Valid {
			task.ErrorMessage = errMsg.String
		} else {
			task.ErrorMessage = ""
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetPendingTasks retrieves tasks that are pending execution
func (r *CrawlTaskRepository) GetPendingTasks(ctx context.Context, limit int) ([]*domain.CrawlTask, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, target_url, site_name, selector_config, frequency, last_crawled_at, next_crawl_at, status, error_message, created_at
		FROM crawl_tasks
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending crawl tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.CrawlTask
	for rows.Next() {
		task := &domain.CrawlTask{}
		var selectorConfigJSON []byte
		var lastCrawledAt, nextCrawlAt sql.NullTime
		var errMsg sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.TargetURL,
			&task.SiteName,
			&selectorConfigJSON,
			&task.Frequency,
			&lastCrawledAt,
			&nextCrawlAt,
			&task.Status,
			&errMsg,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crawl task: %w", err)
		}

		if len(selectorConfigJSON) > 0 {
			if err := json.Unmarshal(selectorConfigJSON, &task.SelectorConfig); err != nil {
				return nil, fmt.Errorf("failed to unmarshal selector config: %w", err)
			}
		}

		if lastCrawledAt.Valid {
			t := lastCrawledAt.Time
			task.LastCrawledAt = &t
		}
		if nextCrawlAt.Valid {
			t := nextCrawlAt.Time
			task.NextCrawlAt = &t
		}
		if errMsg.Valid {
			task.ErrorMessage = errMsg.String
		} else {
			task.ErrorMessage = ""
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
