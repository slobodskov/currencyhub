// Package repository provides PostgreSQL implementation of UserRepository
// Handles database operations for user management
package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

// UserRepo implements UserRepository interface for PostgreSQL
// Provides concrete database operations for user data
type UserRepo struct {
	db     *sqlx.DB
	logger slog.Logger
}

// NewUserService creates user data access service
// Initializes with database connection dependency
func NewUserService(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetSubscribedUsers retrieves all users with auto-subscription enabled
// Returns map of user IDs to their update intervals
func (du *UserRepo) GetSubscribedUsers(ctx context.Context) (map[int64]uint, error) {
	type result struct {
		TelegramID   int64 `db:"telegram_id"`
		SendInterval uint  `db:"send_interval"`
	}

	var users []result
	query := `SELECT telegram_id, send_interval FROM users WHERE auto_subscribe = true`
	err := du.db.SelectContext(ctx, &users, query)
	if err != nil {
		du.logger.Error("Failed to get subscribed users", "error", err)
		return nil, fmt.Errorf("failed to get subscribed users: %w", err)
	}

	resultMap := make(map[int64]uint)
	for _, u := range users {
		resultMap[u.TelegramID] = u.SendInterval
	}

	return resultMap, nil
}

// SetAutoSubscribe enables automatic updates for a user
// Configures update interval for auto-subscription
func (du *UserRepo) SetAutoSubscribe(ctx context.Context, userID int64, interval uint) error {
	query := `INSERT INTO users (telegram_id, auto_subscribe, send_interval)
VALUES ($2, true, $1)
ON CONFLICT (telegram_id) 
DO UPDATE SET 
    auto_subscribe = EXCLUDED.auto_subscribe,
    send_interval = EXCLUDED.send_interval;`
	_, err := du.db.ExecContext(ctx, query, interval, userID)
	if err != nil {
		du.logger.Error("Failed to set auto subscribe", "userID", userID, "error", err)
		return fmt.Errorf("failed to set auto subscribe: %w", err)
	}
	return nil
}

// DisableAutoSubscribe disables automatic updates for a user
// Turns off auto-subscription feature for specified user
func (du *UserRepo) DisableAutoSubscribe(ctx context.Context, userID int64) error {
	query := `UPDATE users SET auto_subscribe = false, send_interval = 0 WHERE telegram_id = $1`
	_, err := du.db.ExecContext(ctx, query, userID)
	if err != nil {
		du.logger.Error("Failed to disable auto subscribe", "userID", userID, "error", err)
		return fmt.Errorf("failed to disable auto subscribe: %w", err)
	}
	return nil
}

// GetUserSendInterval retrieves update interval for a specific user
// Returns configured update frequency in minutes
func (du *UserRepo) GetUserSendInterval(ctx context.Context, userID int64) (uint, error) {
	var interval uint
	query := `SELECT send_interval FROM users WHERE telegram_id = $1`
	err := du.db.GetContext(ctx, &interval, query, userID)
	if err != nil {
		du.logger.Error("Failed to get user send interval", "userID", userID, "error", err)
		return 0, fmt.Errorf("failed to get user send interval: %w", err)
	}
	return interval, nil
}
