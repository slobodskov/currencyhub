// Package database provides database schema migration functionality
// Handles database structure creation and updates
package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
	"strings"
)

// Migrator manages database schema migrations
// Handles application of SQL migration scripts
type Migrator struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewMigrator creates new database migrator instance
// Initializes with database connection and logger
func NewMigrator(db *sqlx.DB, logger *slog.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// Migrate applies database schema migrations
// Executes SQL scripts to create or update database structure
func (m *Migrator) Migrate(ctx context.Context) error {
	sqlBytes, err := os.ReadFile("migrations/migrations_up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations file: %w", err)
	}

	queries := strings.Split(string(sqlBytes), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := m.db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute migration: %w\nQuery: %s", err, query)
		}
	}

	m.logger.Info("Migrations applied successfully")
	return nil
}

// Rollback reverts database schema changes
// Executes rollback scripts to undo migrations
func (m *Migrator) Rollback(ctx context.Context, db *sqlx.DB, logger slog.Logger) error {
	sqlBytes, err := os.ReadFile("migrations/migrations_up.sql")
	if err != nil {
		return fmt.Errorf("failed to read rollback file: %w", err)
	}

	queries := strings.Split(string(sqlBytes), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute rollback: %w\nQuery: %s", err, query)
		}
	}

	logger.Info("Migrations rolled back successfully")
	return nil
}
