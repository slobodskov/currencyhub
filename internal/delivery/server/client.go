package server

import (
	"currencyhub/internal/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// CurrencyHandler manages HTTP requests for currency data
// Contains business logic controller for currency operations
type CurrencyHandler struct {
	currencyUseCase *usecase.CurrencyUseCase
}

// NewCurrencyHandler creates new CurrencyHandler instance
// Initializes with currency use case dependency
func NewCurrencyHandler(currencyUseCase *usecase.CurrencyUseCase) *CurrencyHandler {
	return &CurrencyHandler{currencyUseCase: currencyUseCase}
}

// Routes configures HTTP routes for currency endpoints
// Returns configured router with all currency endpoints
func (h *CurrencyHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(prometheusMiddleware)

	r.Handle("/metrics", promhttp.Handler())

	r.Get("/rates", h.GetRates)
	r.Get("/rates/{currency}", h.GetCurrencyRate)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), // URL к сгенерированному файлу doc.json
	))
	return r
}
