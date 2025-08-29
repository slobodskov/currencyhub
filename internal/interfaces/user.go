// Package repository defines data storage interfaces
// Abstracts database operations for user management
package interfaces

import (
	"context"
	_ "currencyhub/internal/entities"
)

// UserRepository defines interface for user data operations
// Provides contract for database interactions with user preferences
type UserRepository interface {
	GetSubscribedUsers(ctx context.Context) (map[int64]uint, error)          // Gets all users with auto-subscription enabled
	SetAutoSubscribe(ctx context.Context, userID int64, interval uint) error // Enables auto-subscription for user
	DisableAutoSubscribe(ctx context.Context, userID int64) error            // Disables auto-subscription for user
	GetUserSendInterval(ctx context.Context, userID int64) (uint, error)     // Gets user's update interval setting
}
