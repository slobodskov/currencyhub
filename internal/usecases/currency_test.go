package usecase

import (
	"context"
	"currencyhub/internal/entities"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockCurrencyRepository struct {
	mock.Mock
}

func (m *MockCurrencyRepository) GetLatestByCurrency(ctx context.Context, currencyID string) (*entities.CurrencyRate, error) {
	args := m.Called(ctx, currencyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CurrencyRate), args.Error(1)
}

func (m *MockCurrencyRepository) GetRates(ctx context.Context) ([]*entities.CurrencyRate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.CurrencyRate), args.Error(1)
}

func (m *MockCurrencyRepository) WriteToBase(ctx context.Context, rate *entities.CurrencyRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockCurrencyRepository) UpdateHourlyStats(ctx context.Context, rate *entities.CurrencyRate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

func (m *MockCurrencyRepository) GetDataInfo(ctx context.Context, coinID string) (*entities.CurrencyRate, error) {
	args := m.Called(ctx, coinID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CurrencyRate), args.Error(1)
}

func (m *MockCurrencyRepository) UpdateDailyData(rate *entities.CurrencyRate) error {
	args := m.Called(rate)
	return args.Error(0)
}

func (m *MockCurrencyRepository) CheckList(coin string) bool {
	args := m.Called(coin)
	return args.Bool(0)
}

func (m *MockCurrencyRepository) SavePrice(ctx context.Context, coinID string, price float64) error {
	args := m.Called(ctx, coinID, price)
	return args.Error(0)
}

func TestCurrencyUseCase_GetRates(t *testing.T) {
	mockRepo := new(MockCurrencyRepository)
	useCase := NewCurrencyUseCase(mockRepo)

	expectedRates := []*entities.CurrencyRate{
		{CurrencyID: "bitcoin", CurrentPrice: 50000},
		{CurrencyID: "ethereum", CurrentPrice: 3000},
	}

	mockRepo.On("GetRates", mock.Anything).Return(expectedRates, nil)

	rates, err := useCase.GetRates(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedRates, rates)
	mockRepo.AssertExpectations(t)
}

func TestCurrencyUseCase_GetRates_Error(t *testing.T) {
	mockRepo := new(MockCurrencyRepository)
	useCase := NewCurrencyUseCase(mockRepo)

	mockRepo.On("GetRates", mock.Anything).Return(nil, errors.New("database error"))

	rates, err := useCase.GetRates(context.Background())

	assert.Error(t, err)
	assert.Nil(t, rates)
	mockRepo.AssertExpectations(t)
}

func TestCurrencyUseCase_GetCurrencyRate(t *testing.T) {
	mockRepo := new(MockCurrencyRepository)
	useCase := NewCurrencyUseCase(mockRepo)

	expectedRate := &entities.CurrencyRate{
		CurrencyID:    "bitcoin",
		CurrentPrice:  50000,
		ChangePercent: 2.5,
	}

	mockRepo.On("GetLatestByCurrency", mock.Anything, "bitcoin").Return(expectedRate, nil)

	rate, err := useCase.GetLatestByCurrency(context.Background(), "bitcoin")

	assert.NoError(t, err)
	assert.Equal(t, expectedRate, rate)
	mockRepo.AssertExpectations(t)
}

func TestCurrencyUseCase_GetCurrencyRate_DBError(t *testing.T) {
	mockRepo := new(MockCurrencyRepository)
	useCase := NewCurrencyUseCase(mockRepo)

	mockRepo.On("GetLatestByCurrency", mock.Anything, "bitcoin").Return(nil, errors.New("database error"))

	rate, err := useCase.GetLatestByCurrency(context.Background(), "bitcoin")

	assert.Error(t, err)
	assert.Nil(t, rate)
	mockRepo.AssertExpectations(t)
}

func TestCurrencyUseCase_CheckList(t *testing.T) {
	mockRepo := new(MockCurrencyRepository)
	useCase := NewCurrencyUseCase(mockRepo)

	mockRepo.On("CheckList", "bitcoin").Return(true)
	mockRepo.On("CheckList", "invalid").Return(false)

	assert.True(t, useCase.CheckList("bitcoin"))
	assert.False(t, useCase.CheckList("invalid"))
	mockRepo.AssertExpectations(t)
}
