package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetSubscribedUsers(ctx context.Context) (map[int64]uint, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[int64]uint), args.Error(1)
}

func (m *MockUserRepository) SetAutoSubscribe(ctx context.Context, userID int64, interval uint) error {
	args := m.Called(ctx, userID, interval)
	return args.Error(0)
}

func (m *MockUserRepository) DisableAutoSubscribe(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserSendInterval(ctx context.Context, userID int64) (uint, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(uint), args.Error(1)
}

func TestUserUseCase_SetAutoSubscribe(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.On("SetAutoSubscribe", mock.Anything, int64(123), uint(10)).Return(nil)

	err := useCase.SetAutoSubscribe(context.Background(), 123, 10)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DisableAutoSubscribe(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)

	mockRepo.On("DisableAutoSubscribe", mock.Anything, int64(123)).Return(nil)

	err := useCase.DisableAutoSubscribe(context.Background(), 123)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetSubscribedUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)

	expectedUsers := map[int64]uint{123: 10, 456: 15}
	mockRepo.On("GetSubscribedUsers", mock.Anything).Return(expectedUsers, nil)

	users, err := useCase.GetSubscribedUsers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}
