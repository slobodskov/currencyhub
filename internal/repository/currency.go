// Package repository provides PostgreSQL implementation of CurrencyRepository
// Handles actual database operations for currency data management
package repository

import (
	"context"
	"currencyhub/internal/entities"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

// CurrencyRepo implements CurrencyRepository interface for PostgreSQL
// Provides concrete database operations for currency data
type CurrencyRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewCurrencyRepo creates new currency repository instance
// Initializes with database connection dependency
func NewCurrencyRepo(db *sqlx.DB) *CurrencyRepo {
	return &CurrencyRepo{db: db}
}

// GetLatestByCurrency retrieves latest rate for specific currency
// Returns most recent currency rate data from database
func (r *CurrencyRepo) GetLatestByCurrency(ctx context.Context, currencyID string) (*entities.CurrencyRate, error) {
	if !r.CheckList(currencyID) {
		return nil, fmt.Errorf("currency not found: %s", currencyID)
	}

	var cr entities.CurrencyRate
	query := `SELECT currency_id, current_price, min_price, max_price, change_percent,
		hour_min_price, hour_max_price, time_stamp, date
		FROM currencies WHERE currency_id = $1 ORDER BY time_stamp DESC LIMIT 1`

	err := r.db.GetContext(ctx, &cr, query, currencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("currency not found: %s", currencyID)
		}
		r.logger.Error("Failed to get currency rate", "currency", currencyID, "error", err)
		return nil, err
	}
	return &cr, nil
}

// GetRates retrieves latest rates for all available cryptocurrencies
// Returns current market data for all supported currencies
func (r *CurrencyRepo) GetRates(ctx context.Context) ([]*entities.CurrencyRate, error) {
	query := `
        SELECT DISTINCT ON (currency_id) 
            currency_id, current_price, min_price, max_price, change_percent,
            hour_min_price, hour_max_price, time_stamp, date
        FROM currencies 
        ORDER BY currency_id, time_stamp DESC
    `

	var crs []*entities.CurrencyRate
	err := r.db.SelectContext(ctx, &crs, query)
	if err != nil {
		r.logger.Error("Failed to get rates", "error", err)
		return nil, err
	}

	return crs, nil
}

// CheckList validates currency exists in supported list
// Returns true if currency is supported by application
func (r *CurrencyRepo) CheckList(coin string) bool {

	for _, i := range entities.CurrencyList {
		if i == coin {
			return true
		}
	}
	return false
}

// SavePrice updates currency data in database with new price information
// Maintains hourly and daily statistics for price tracking
func (r *CurrencyRepo) SavePrice(ctx context.Context, coinID string, price float64) error {
	existingRate, err := r.GetDataInfo(ctx, coinID)
	if err != nil {
		if err == sql.ErrNoRows {
			now := time.Now().UTC()
			existingRate = &entities.CurrencyRate{
				CurrencyID:   coinID,
				CurrentPrice: price,
				MinPrice:     price,
				MaxPrice:     price,
				HourMinPrice: price,
				HourMaxPrice: price,
				TimeStamp:    now,
				Date:         time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			}
		} else {
			return fmt.Errorf("failed to get currency data for %s: %w", coinID, err)
		}
	} else {
		existingRate.CurrentPrice = price
	}

	if err := r.UpdateDailyData(existingRate); err != nil {
		return fmt.Errorf("failed to update daily stats for %s: %w", coinID, err)
	}

	if err := r.UpdateHourlyStats(existingRate); err != nil {
		return fmt.Errorf("failed to update hourly stats for %s: %w", coinID, err)
	}

	if err := r.WriteToBase(ctx, existingRate); err != nil {
		return fmt.Errorf("failed to save currency data for %s: %w", coinID, err)
	}

	return nil
}

// GetDataInfo retrieves currency information from database
// Returns complete currency rate record for specified coin
// Only for SavePrice
func (r *CurrencyRepo) GetDataInfo(ctx context.Context, coinID string) (*entities.CurrencyRate, error) {
	var currency entities.CurrencyRate

	query := `SELECT currency_id, current_price, min_price, max_price, change_percent, 
		hour_min_price, hour_max_price, time_stamp, date 
		FROM currencies WHERE currency_id = $1`

	err := r.db.QueryRowContext(ctx, query, coinID).Scan(
		&currency.CurrencyID,
		&currency.CurrentPrice,
		&currency.MinPrice,
		&currency.MaxPrice,
		&currency.ChangePercent,
		&currency.HourMinPrice,
		&currency.HourMaxPrice,
		&currency.TimeStamp,
		&currency.Date,
	)
	if err != nil {
		return nil, err
	}
	return &currency, nil
}

// UpdateDailyData handles daily price statistics updates.
// Maintains daily min/max price tracking.
// Only for SavePrice
func (r *CurrencyRepo) UpdateDailyData(rate *entities.CurrencyRate) error {
	now := time.Now()
	currentDate := now.UTC().Truncate(24 * time.Hour)

	if rate.Date.Before(currentDate) {
		rate.MinPrice = rate.CurrentPrice
		rate.MaxPrice = rate.CurrentPrice
		rate.Date = currentDate
	} else {
		if rate.CurrentPrice < rate.MinPrice {
			rate.MinPrice = rate.CurrentPrice
		}
		if rate.CurrentPrice > rate.MaxPrice {
			rate.MaxPrice = rate.CurrentPrice
		}
	}

	return nil
}

// UpdateHourlyStats handles hourly price tracking and statistics
// Calculates hourly min/max prices and change percentages
// Only for SavePrice
func (r *CurrencyRepo) UpdateHourlyStats(rate *entities.CurrencyRate) error {
	if time.Since(rate.TimeStamp) >= time.Hour {
		rate.HourMinPrice = rate.CurrentPrice
		rate.HourMaxPrice = rate.CurrentPrice
		rate.TimeStamp = time.Now()
	} else {
		if rate.CurrentPrice < rate.HourMinPrice {
			rate.HourMinPrice = rate.CurrentPrice
		}
		if rate.CurrentPrice > rate.HourMaxPrice {
			rate.HourMaxPrice = rate.CurrentPrice
		}
	}

	if rate.HourMinPrice == 0 {
		rate.ChangePercent = 0
	} else {
		rate.ChangePercent = (rate.HourMaxPrice - rate.HourMinPrice) / rate.HourMinPrice * 100
	}
	return nil
}

// WriteToBase inserts or updates currency record with statistics
// Implements upsert operation for currency data management
// Only for SavePrice
func (r *CurrencyRepo) WriteToBase(ctx context.Context, currency *entities.CurrencyRate) error {
	query := `
        INSERT INTO currencies 
            (currency_id, current_price, min_price, max_price, change_percent, 
            hour_min_price, hour_max_price, time_stamp, date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (currency_id) 
        DO UPDATE SET
            current_price = EXCLUDED.current_price,
            min_price = LEAST(currencies.min_price, EXCLUDED.min_price),
            max_price = GREATEST(currencies.max_price, EXCLUDED.max_price),
            change_percent = EXCLUDED.change_percent,
            hour_min_price = LEAST(currencies.hour_min_price, EXCLUDED.hour_min_price),
            hour_max_price = GREATEST(currencies.hour_max_price, EXCLUDED.hour_max_price),
            time_stamp = EXCLUDED.time_stamp,
            date = EXCLUDED.date
    `
	_, err := r.db.ExecContext(
		ctx,
		query,
		currency.CurrencyID,
		currency.CurrentPrice,
		currency.MinPrice, // min_price
		currency.MaxPrice, // max_price
		currency.ChangePercent,
		currency.HourMinPrice, // hour_min_price
		currency.HourMaxPrice, // hour_max_price
		currency.TimeStamp,    // time_stamp
		currency.Date,         // date
	)

	if err != nil {
		return fmt.Errorf("failed to update currency data: %w", err)
	}
	return nil
}
