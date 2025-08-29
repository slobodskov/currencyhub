// Currency use cases interfaces.
// Defines contracts for:
// - Currency data retrieval
// - Statistics calculation
// - Business logic separation
package usecase

import (
	"context"
	"currencyhub/internal/entities"
	"currencyhub/internal/interfaces"
)

// CurrencyUseCase provides business logic operations for currency data
// Acts as an intermediary between delivery layer (handlers) and repository layer
type CurrencyUseCase struct {
	currencyRepo interfaces.CurrencyRepository
}

// NewCurrencyUseCase creates a new instance of CurrencyUseCase
// with the provided currency repository dependency
func NewCurrencyUseCase(currencyRepo interfaces.CurrencyRepository) *CurrencyUseCase {
	return &CurrencyUseCase{currencyRepo: currencyRepo}
}

// GetRates retrieves latest rates for all available cryptocurrencies
// Returns a slice of CurrencyRate entities or error if operation fails
func (uc *CurrencyUseCase) GetRates(ctx context.Context) ([]*entities.CurrencyRate, error) {
	return uc.currencyRepo.GetRates(ctx)
}

// GetLatestByCurrency retrieves most recent rate for specific cryptocurrency
func (uc *CurrencyUseCase) GetLatestByCurrency(ctx context.Context, currencyID string) (*entities.CurrencyRate, error) {
	return uc.currencyRepo.GetLatestByCurrency(ctx, currencyID)
}

// CheckList validates currency exists in supported list
// Returns true if currency is supported by application
func (uc *CurrencyUseCase) CheckList(coin string) bool {
	return uc.currencyRepo.CheckList(coin)
}

// SavePrice updates currency data in database with new price information
// Maintains hourly and daily statistics for price tracking
func (uc *CurrencyUseCase) SavePrice(ctx context.Context, coinID string, price float64) error {
	return uc.currencyRepo.SavePrice(ctx, coinID, price)
}
