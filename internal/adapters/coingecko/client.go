// Package coingecko provides CoinGecko API client implementation
// Handles cryptocurrency price data fetching and storage
package coingecko

import (
	"context"
	"currencyhub/internal/usecases"
	"log/slog"
	"net/http"
	"time"
)

// Client manages interactions with CoinGecko API
// Handles HTTP requests and response processing for cryptocurrency data
type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
	apiKey     string
	repo       *usecase.CurrencyUseCase
}

// NewClient creates CoinGecko API client instance
// Initializes with configured HTTP client and dependencies
func NewClient(logger *slog.Logger, apiKey string, repo *usecase.CurrencyUseCase) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
		apiKey:     apiKey,
		repo:       repo,
	}
}

// Run starts periodic currency update scheduler
// Executes updates at fixed intervals with graceful shutdown
func (c *Client) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	update := func() {
		c.logger.Info("Starting currency update")
		if err := c.GetPrices(ctx); err != nil {
			c.logger.Error("Failed to update prices", "error", err)
			return
		}
		c.logger.Info("Currency update completed successfully")
	}

	update()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Updater stopped")
			return
		case <-ticker.C:
			update()
		}
	}
}
