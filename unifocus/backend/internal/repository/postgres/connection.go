package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/unifocus/backend/internal/config"
	"github.com/unifocus/backend/pkg/logger"
)

// DB wraps sql.DB with additional functionality
type DB struct {
	*sql.DB
}

// NewDatabase creates a new PostgreSQL database connection with connection pool
// It configures connection pool parameters based on the provided config
// and performs a health check (Ping) to ensure the connection is valid
func NewDatabase(cfg *config.DatabaseConfig) (*DB, error) {
	// Open database connection
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Perform health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Infof("Database connection established: %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// HealthCheck performs a health check on the database connection
func (d *DB) HealthCheck(ctx context.Context) error {
	return d.PingContext(ctx)
}

// Transaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (d *DB) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Errorf("failed to rollback transaction: %v", rbErr)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				logger.Errorf("failed to commit transaction: %v", err)
			}
		}
	}()

	err = fn(tx)
	return err
}
