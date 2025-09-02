// app contains main application initialization and orchestration logic
// Coordinates between different application components and services
package app

import (
	"context"
	"currencyhub/config"
	"currencyhub/internal/adapters/coingecko"
	"currencyhub/internal/adapters/postgres"
	"currencyhub/internal/delivery/server"
	"currencyhub/internal/delivery/telegram"
	"currencyhub/internal/infrastructure/database"
	log "currencyhub/internal/infrastructure/logger"
	"currencyhub/internal/infrastructure/shutdown"
	"currencyhub/internal/repository"
	"currencyhub/internal/usecases"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

// Run starts all application services and components
// Initializes use cases, starts scheduler, bot and HTTP server
func Run() error {

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := log.SetupLogger(cfg.Logging.File)

	dsn := database.LoadConfig(cfg)

	var (
		db   *sqlx.DB
		err1 error
	)
	for i := 0; i < 5; i++ {
		db, err1 = database.NewPostgresDB(dsn)
		if err1 == nil {
			break
		}
		logger.Warn("Failed to connect to database, retrying...", "attempt", i+1)
		time.Sleep(5 * time.Second)
	}

	migrator := postgres.NewMigrator(db, logger)
	if err := migrator.Migrate(context.Background()); err != nil {
		return err
	}

	currencyRepo := repository.NewCurrencyRepo(db)
	userRepo := repository.NewUserService(db)

	currencyService := usecase.NewCurrencyUseCase(currencyRepo)
	userService := usecase.NewUserUseCase(userRepo)

	receiver := coingecko.NewClient(logger, cfg.Coingecko.APIKey, currencyService)
	go receiver.Run(ctx)

	bot, err := telegram.NewBot(userService, currencyService, logger, cfg.Telegram.Token)
	if err != nil {
		return fmt.Errorf("telegram bot not created: %w", err)
	}
	go bot.Run(ctx)

	handler := server.NewCurrencyHandler(currencyService)
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: handler.Routes(),
	}

	go func() {
		logger.Info("Starting HTTP server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", "error", err)
		}
	}()

	rollbackHook := func() error {
		logger.Info("Executing database rollback during shutdown")
		return migrator.Rollback(context.Background(), db, *logger)
	}

	return shutdown.WaitForShutdown(ctx, server, db, logger, rollbackHook)
}
