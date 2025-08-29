package database

import (
	"currencyhub/config"
	"fmt"
)

// LoadConfig constructs PostgreSQL connection string from configuration
// Returns formatted Data Source Name for database connection
func LoadConfig(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
}
