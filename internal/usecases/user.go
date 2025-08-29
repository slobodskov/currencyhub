// User management use cases.
// Contains:
// - User entities with business logic
// - Interface for user operations
// - Subscription management contracts
package usecase

import (
	"context"
	"currencyhub/internal/interfaces"
)

// UserUseCase struct represents user entities with business logic
type UserUseCase struct {
	userRepo interfaces.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase
// with the provided user repository dependency
func NewUserUseCase(userRepo interfaces.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// GetSubscribedUsers retrieves all users with auto-subscription enabled
func (uc *UserUseCase) GetSubscribedUsers(ctx context.Context) (map[int64]uint, error) {
	return uc.userRepo.GetSubscribedUsers(ctx)
}

// SetAutoSubscribe enables automatic updates for a user
func (uc *UserUseCase) SetAutoSubscribe(ctx context.Context, userID int64, interval uint) error {
	return uc.userRepo.SetAutoSubscribe(ctx, userID, interval)
}

// DisableAutoSubscribe disables automatic updates for a user
func (uc *UserUseCase) DisableAutoSubscribe(ctx context.Context, userID int64) error {
	return uc.userRepo.DisableAutoSubscribe(ctx, userID)
}

// GetUserSendInterval returns autosending time in minutes for subscribed users
func (uc *UserUseCase) GetUserSendInterval(ctx context.Context, userID int64) (uint, error) {
	return uc.userRepo.GetUserSendInterval(ctx, userID)
}
