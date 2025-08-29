// Package server handles HTTP request processing for currency data
// Provides REST API endpoints for currency rate information
package server

import (
	"currencyhub/internal/entities"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

// GetRates handles HTTP GET request for all currency rates
// @Summary Получить все курсы валют
// @Description Возвращает список всех доступных курсов криптовалют
// @Tags rates
// @Produce plain
// @Success 200 {string} string "Курсы валют"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /rates [get]
func (h *CurrencyHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rates, err := h.currencyUseCase.GetRates(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formattedRates := make([]string, 0, len(rates))
	for _, rate := range rates {
		formattedRate, err := h.FormatOutput(rate)
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			continue
		}
		formattedRates = append(formattedRates, formattedRate)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response := strings.Join(formattedRates, "\r\n\r\n")
	w.Write([]byte(response))
}

// GetCurrencyRate handles HTTP GET request for specific currency rate
// @Summary Получить курс конкретной валюты
// @Description Возвращает детальную информацию по конкретной криптовалюте
// @Tags rates
// @Produce plain
// @Param currency path string true "Идентификатор валюты"
// @Success 200 {string} string "Данные по валюте"
// @Failure 404 {string} string "Валюта не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /rates/{currency} [get]
func (h *CurrencyHandler) GetCurrencyRate(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "currency")

	isExist := h.currencyUseCase.CheckList(currencyID)
	if !isExist {
		http.Error(w, "not found in the list", http.StatusNotFound)
		return
	}

	ctx := r.Context()
	rate, err := h.currencyUseCase.GetLatestByCurrency(ctx, currencyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	formattedRate, err := h.FormatOutput(rate)
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(formattedRate))
}

// FormatOutput formats currency rate data for display
// Returns formatted string with currency information
func (h *CurrencyHandler) FormatOutput(rate *entities.CurrencyRate) (string, error) {
	formatted := fmt.Sprintf(
		"CurrencyID: %s\r\nCurrentPrice: %.2f\r\nMinPrice: %.2f\r\nMaxPrice: %.2f\r\nChangePercent: %.2f%%",
		rate.CurrencyID,
		rate.CurrentPrice,
		rate.MinPrice,
		rate.MaxPrice,
		rate.ChangePercent,
	)
	return formatted, nil
}
