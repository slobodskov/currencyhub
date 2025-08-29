// Package telegram handles Telegram bot functionality
// Processes bot commands and manages user interactions
package telegram

import (
	"context"
	"currencyhub/internal/entities"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

// handleStart processes /start command - welcomes user and shows available commands
func (b *Bot) handleStart(message *tgbotapi.Message) {

	msg := `🤖 💰 Добро пожаловать в Currency Hub Bot! 

📋 Доступные команды:
/rates - показать все курсы валют 📊
/rates [валюта] - показать курс конкретной валюты 📈
/coins - список всех доступных валют 🪙
/start_auto [мин] - запустить автоподписку 🔔
/stop_auto - остановить автоподписку 🔕
/help - показать справку ❓`

	b.sendMessage(message.Chat.ID, msg)
}

// HandleRates processes /rates command - shows currency rates (all or specific)
func (b *Bot) handleRates(ctx context.Context, message *tgbotapi.Message) {
	args := strings.Split(message.Text, " ")
	if len(args) > 1 {
		currencyID := strings.ToLower(args[1])
		isExist := b.currencyUseCase.CheckList(currencyID)
		if !isExist {
			b.sendMessage(message.Chat.ID, "❌ Валюта не найдена")
			return
		}
		rate, err := b.currencyUseCase.GetLatestByCurrency(ctx, currencyID)
		if err != nil {
			b.sendMessage(message.Chat.ID, "❌ Ошибка при получении, попробуйте позже")
			return
		}

		trendEmoji := "➡️"
		if rate.ChangePercent > 0 {
			trendEmoji = "📈"
		} else if rate.ChangePercent < 0 {
			trendEmoji = "📉"
		}

		msg := fmt.Sprintf("💰 Курс %s:\n📊 Текущий: $%.2f\n📉 Мин. за день: $%.2f\n📈 Макс. за день: $%.2f\n%s Изменение за час: %.2f%%",
			currencyID, rate.CurrentPrice, rate.MinPrice, rate.MaxPrice, trendEmoji, rate.ChangePercent)
		b.sendMessage(message.Chat.ID, msg)
	} else {
		rates, err := b.currencyUseCase.GetRates(ctx)
		if err != nil {
			b.logger.Error("Failed to get rates", "error", err)
			b.sendMessage(message.Chat.ID, "❌ Ошибка получения курсов")
			return
		}

		var msg strings.Builder
		msg.WriteString("📊 Текущие курсы:\n\n")
		for _, rate := range rates {
			trendEmoji := "➡️"
			if rate.ChangePercent > 0 {
				trendEmoji = "📈"
			} else if rate.ChangePercent < 0 {
				trendEmoji = "📉"
			}

			msg.WriteString(fmt.Sprintf("💰 %s: $%.2f %s(%.2f%%)\n",
				rate.CurrencyID, rate.CurrentPrice, trendEmoji, rate.ChangePercent))
		}
		b.sendMessage(message.Chat.ID, msg.String())
	}
}

// HandleCoins processes /coins command - shows available cryptocurrencies
func (b *Bot) handleCoins(message *tgbotapi.Message) {

	coins := entities.CurrencyList

	msg := "📋 Доступные валюты:\n" + strings.Join(coins, "\n")
	b.sendMessage(message.Chat.ID, msg)
}

// HandleStartAuto processes /start_auto command - enables automatic updates with time choice option
func (b *Bot) handleStartAuto(ctx context.Context, message *tgbotapi.Message) {
	args := strings.Split(message.Text, " ")
	interval := uint(10)

	if len(args) > 1 {
		if minutes, err := strconv.Atoi(args[1]); err == nil && minutes > 0 {
			interval = uint(minutes)
		} else if minutes, err := strconv.Atoi(args[1]); err == nil && minutes < 1 {
			b.sendMessage(message.Chat.ID, "❌ Неверный интервал. Используйте число больше 0")
			return
		}
	}

	err := b.userUseCase.SetAutoSubscribe(ctx, message.Chat.ID, interval)
	if err != nil {
		b.logger.Error("Failed to set auto subscribe", "error", err)
		b.sendMessage(message.Chat.ID, "❌ Ошибка включения автоподписки")
		return
	}

	msg := fmt.Sprintf("🔔 Автоподписка включена. Частота обновлений (минуты): ", interval)
	b.sendMessage(message.Chat.ID, msg)

}

// HandleStopAuto processes /stop_auto command - disables automatic updates
func (b *Bot) handleStopAuto(ctx context.Context, message *tgbotapi.Message) {
	err := b.userUseCase.DisableAutoSubscribe(ctx, message.Chat.ID)
	if err != nil {
		b.logger.Error("Failed to disable auto subscribe", "error", err)
		b.sendMessage(message.Chat.ID, "❌ Ошибка отключения автоподписки")
		return
	}

	b.sendMessage(message.Chat.ID, "🔕 Автоподписка отключена")
}

// handleHelp processes /help command - shows available commands
func (b *Bot) handleHelp(message *tgbotapi.Message) {
	msg := `🤖 💰 Доступные команды:

📊 /rates - показать все курсы валют
📈 /rates [валюта] - показать курс конкретной валюты
🪙 /coins - список всех доступных валют
🔔 /start_auto [мин] - запустить автоподписку
🔕 /stop_auto - остановить автоподписку
❓ /help - показать справку`

	b.sendMessage(message.Chat.ID, msg)
}

// sendMessage sends text message to Telegram chat
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.Api.Send(msg)
	if err != nil {
		b.logger.Error("Failed to send message", "error", err)
	}
}

// sendUpdates periodically sends currency updates to subscribed users
func (b *Bot) sendUpdates(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.logger.Info("sending starts")
			b.sendCurrencyUpdates(ctx)
		}
	}
}

// sendCurrencyUpdates gets user with autoupdate setting and sends stats
func (b *Bot) sendCurrencyUpdates(ctx context.Context) {
	users, err := b.userUseCase.GetSubscribedUsers(ctx)
	if err != nil {
		b.logger.Error("Failed to get subscribed users", "error", err)
		return
	} else {
		b.logger.Info("got subscribed users")
	}

	rates, err := b.currencyUseCase.GetRates(ctx)
	if err != nil {
		b.logger.Error("Failed to get currency rates", "error", err)
		return
	} else {
		b.logger.Info("got currency rates")
	}

	var msg strings.Builder
	msg.WriteString("🔔 Автообновление курсов:\n\n")
	for _, rate := range rates {
		trendEmoji := "➡️"
		if rate.ChangePercent > 0 {
			trendEmoji = "📈"
		} else if rate.ChangePercent < 0 {
			trendEmoji = "📉"
		}

		msg.WriteString(fmt.Sprintf("💰 %s: $%.2f %s(%.2f%%)\n",
			rate.CurrencyID, rate.CurrentPrice, trendEmoji, rate.ChangePercent))
	}
	message := msg.String()

	now := time.Now()
	for userID, interval := range users {
		lastSent, exists := b.lastSentMap[userID]
		shouldSend := !exists || now.Sub(lastSent) >= time.Duration(interval)*time.Minute

		if shouldSend {
			b.logger.Info("sending update")
			b.sendMessage(userID, message)
			b.lastSentMap[userID] = now
		}
	}
}
