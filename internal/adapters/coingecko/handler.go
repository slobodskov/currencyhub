package coingecko

import (
	"context"
	"currencyhub/internal/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetPrices fetches cryptocurrency prices from CoinGecko API
// Retrieves current market data for all supported currencies
func (c *Client) GetPrices(ctx context.Context) error {
	coinIDs := entities.CurrencyList

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd",
		strings.Join(coinIDs, ","))

	if c.apiKey != "" {
		url += "&x_cg_demo_api_key=" + c.apiKey
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.Error("Failed to create request", "error", err)
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to get prices", "error", err)
		return err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.logger.Error("Failed to decode response", "error", err)
		return err
	}

	for coinID, priceData := range data {
		price, ok := priceData["usd"]
		if !ok {
			c.logger.Warn("Price not found for coin", "coin", coinID)
			continue
		}

		if err := c.repo.SavePrice(ctx, coinID, price); err != nil {
			c.logger.Error("Failed to save price", "coin", coinID, "error", err)
		}
	}

	return nil
}
