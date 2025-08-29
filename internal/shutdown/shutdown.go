package shutdown

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownHook receives extra arguments for WaitForShutdown
type ShutdownHook func() error

// WaitForShutdown handles graceful application shutdown
// Listens for termination signals and cleans up resources
func WaitForShutdown(ctx context.Context, server *http.Server, db *sqlx.DB, logger *slog.Logger, hooks ...ShutdownHook) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received")

	for _, hook := range hooks {
		if hook != nil {
			if err := hook(); err != nil {
				logger.Error("Shutdown hook failed", "error", err)
			}
		}
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
