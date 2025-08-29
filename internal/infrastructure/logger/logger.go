package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// SetupLogger configures application logger with specified output
// Creates file logger or stdout logger based on configuration
func SetupLogger(logFile string) *slog.Logger {
	var logHandler slog.Handler

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("Failed to open log file: %w", err)
			return nil
		}
		logHandler = slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(logHandler)
}
