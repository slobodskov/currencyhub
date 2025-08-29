// Package entities defines domain entities and business objects
// Contains core data structures used throughout the application

// Domain entities for cryptocurrency rate data.
// CurrencyRate struct contains:
// - Currency identification and current price
// - Daily minimum and maximum prices
// - Hourly price extremes
// - Price change percentages
// - Timestamp and date information
package entities

import "time"

// CurrencyRate represents cryptocurrency rate information
// Stores current and historical price data with statistics
type CurrencyRate struct {
	CurrencyID    string    `db:"currency_id"`    // Unique cryptocurrency identifier
	CurrentPrice  float64   `db:"current_price"`  // Current market price in USD
	MinPrice      float64   `db:"min_price"`      // Daily minimum price
	MaxPrice      float64   `db:"max_price"`      // Daily maximum price
	ChangePercent float64   `db:"change_percent"` // Hourly price change percentage
	HourMinPrice  float64   `db:"hour_min_price"` // Hourly minimum price
	HourMaxPrice  float64   `db:"hour_max_price"` // Hourly maximum price
	TimeStamp     time.Time `db:"time_stamp"`     // Last update timestamp
	Date          time.Time `db:"date"`           // Date for daily statistics
}
