# Currency Hub

A cryptocurrency tracking application that provides real-time currency rates, automated updates, and Telegram bot integration for convenient access to cryptocurrency market data.

## Features

- **Real-time Cryptocurrency Data**: Fetches current prices from CoinGecko API
- **Telegram Bot Integration**: Interactive bot with commands for currency information
- **Automated Updates**: Scheduled price updates and user notifications
- **REST API**: HTTP endpoints for accessing currency data
- **PostgreSQL Storage**: Persistent data storage with proper schema management
- **Swagger Documentation**: API documentation automatically generated

## Supported Cryptocurrencies

- Bitcoin, Ethereum, Tether, Binance Coin, Solana
- USD Coin, Ripple, The Open Network, Dogecoin, Cardano
- Shiba Inu, Avalanche, Polkadot, Tron, Chainlink
- Polygon, Bitcoin Cash, Litecoin, Uniswap, DAI

## Prerequisites

- Go 1.24.5 or higher
- Docker and Docker Compose
- PostgreSQL 15
- Telegram Bot Token ([Create one here](https://t.me/BotFather))
- CoinGecko API Key (Optional - for higher rate limits)

## Installation

1. **Clone the repository**
   ```bash
   git clone <your-gitlab-repo-url>
   cd currencyhub1
2. **Set up environment variables**
   cp .env.example .env
3. **Edit .env with your actual values**
   TELEGRAM_TOKEN=your_telegram_bot_token
   COINGECKO_API_KEY=your_coingecko_api_key
   DB_PASSWORD=password
4. **Build and run with Docker**
   docker-compose up --build
5. **Or run locally**
   # Install dependencies
    go mod download

    # Run database
    docker-compose up postgres -d

    # Run application
    go run main.go

## API Endpoints
   **REST API**
   - GET /rates - Get all currency rates

   - GET /rates/{currency} - Get specific currency rate

   - GET /swagger/ - Swagger API documentation

   **Telegram Bot Commands**
    
 /start - Welcome message and command list

/rates - Show all currency rates

/rates [currency] - Show specific currency rate

/coins - List available cryptocurrencies

/start_auto [min] - Enable auto-updates (default: 10 min)

/stop_auto - Disable auto-updates

/help - Show help information

**Health Check**

curl http://localhost:8080/rates