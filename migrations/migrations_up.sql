
CREATE TABLE IF NOT EXISTS currencies (
                                          currency_id TEXT PRIMARY KEY,
                                          current_price DECIMAL NOT NULL,
                                          min_price DECIMAL NOT NULL,
                                          max_price DECIMAL NOT NULL,
                                          change_percent DECIMAL NOT NULL,
                                          hour_min_price DECIMAL NOT NULL,
                                          hour_max_price DECIMAL NOT NULL,
                                          time_stamp TIMESTAMP NOT NULL,
                                          date DATE NOT NULL
);

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
                                     telegram_id BIGINT PRIMARY KEY,
                                     auto_subscribe BOOLEAN DEFAULT FALSE,
                                     send_interval INTEGER DEFAULT 10
);

-- Создание индексов для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_currencies_timestamp ON currencies(time_stamp);
CREATE INDEX IF NOT EXISTS idx_currencies_currency_id ON currencies(currency_id);