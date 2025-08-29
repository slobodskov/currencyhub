// Package entities defines domain entities and business objects
// Contains core data structures used throughout the application

// Package domain contains core business entities.
// User struct represents a Telegram bot user with:
// - TelegramID: Unique user identifier from Telegram
// - AutoSubscribe: Flag for automatic currency updates subscription
// - SendInterval: Frequency of updates in minutes
package entities

// User represents Telegram bot user with preferences
// Stores user settings for notification preferences
type User struct {
	TelegramID    int64 `db:"telegram_id"`    // Unique Telegram user identifier
	AutoSubscribe bool  `db:"auto_subscribe"` // Automatic updates subscription status
	SendInterval  uint  `db:"send_interval"`  // Update frequency in minutes
}
