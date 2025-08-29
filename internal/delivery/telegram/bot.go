// Package telegram provides Telegram bot implementation
// Handles bot initialization and message routing
package telegram

import (
	"context"
	"currencyhub/internal/usecases"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"time"
)

// Bot manages Telegram bot lifecycle and message processing
// Handles incoming updates and command routing
type Bot struct {
	logger          *slog.Logger
	userUseCase     *usecase.UserUseCase
	currencyUseCase *usecase.CurrencyUseCase
	lastSentMap     map[int64]time.Time
	Api             *tgbotapi.BotAPI
}

// NewBot creates new Telegram bot instance
// Initializes with message handler and logger
func NewBot(userUseCase *usecase.UserUseCase, currencyUseCase *usecase.CurrencyUseCase, logger *slog.Logger, token string) (*Bot, error) {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{
		userUseCase:     userUseCase,
		currencyUseCase: currencyUseCase,
		lastSentMap:     map[int64]time.Time{},
		Api:             bot,
		logger:          logger,
	}, nil
}

// Run starts Telegram bot update listener
// Processes incoming messages and maintains connection
func (b *Bot) Run(ctx context.Context) {
	b.logger.Info("Starting Telegram bot")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Api.GetUpdatesChan(u)

	go b.sendUpdates(ctx)

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("Telegram bot stopped")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			b.handleMessage(ctx, update.Message)
		}
	}
}

// handleMessage processes incoming Telegram messages and routes to appropriate handlers
// Implements command pattern for bot functionality
func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	if !message.IsCommand() {
		return
	}

	switch message.Command() {
	case "start":
		b.handleStart(ctx, message)
	case "rates":
		b.handleRates(ctx, message)
	case "coins":
		b.handleCoins(ctx, message)
	case "start_auto":
		b.handleStartAuto(ctx, message)
	case "stop_auto":
		b.handleStopAuto(ctx, message)
	case "help":
		b.handleHelp(message)
	default:
		b.sendMessage(message.Chat.ID, "❌ Неизвестная команда. Выберите команду из доступных ниже.")
		b.handleHelp(message)
	}
}
