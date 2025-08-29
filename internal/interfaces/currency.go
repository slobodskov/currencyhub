// Package repository defines data storage interfaces
// Abstracts database operations for currency data management
package interfaces

import (
	"context"
	"currencyhub/internal/entities"
)

// CurrencyRepository defines interface for currency data operations
// Provides contract for database interactions with currency rates
type CurrencyRepository interface {
	GetLatestByCurrency(ctx context.Context, currencyID string) (*entities.CurrencyRate, error) // Gets latest rate for specific currency
	GetRates(ctx context.Context) ([]*entities.CurrencyRate, error)                             // Gets all current currency rates
	CheckList(coin string) bool                                                                 // Validates currency exists in supported list
	SavePrice(ctx context.Context, coinID string, price float64) error                          //updates currency data in database with new price information
}
