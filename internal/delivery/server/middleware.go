package server

import (
	"currencyhub/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		timer := prometheus.NewTimer(monitoring.RequestDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		next.ServeHTTP(wrappedWriter, r)

		monitoring.RequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(wrappedWriter.statusCode),
		).Inc()
	})
}
