# Telegram Raindrop Bot

A Telegram bot that allows users to save links directly to their Raindrop.io account.

<img width="1120" height="1078" alt="image" src="https://github.com/user-attachments/assets/5e129306-8dba-4889-829c-2167c7e4f6f9" />

## Features

- Save links from Telegram to Raindrop.io
- OAuth2 authentication with Raindrop.io
- PostgreSQL database for user management
- Graceful shutdown handling
- Structured logging
- Proper error handling and timeouts

## Architecture

The application follows a clean architecture pattern with the following components:

1. **Main Application** (`internal/app/app.go`) - Orchestrates the entire application lifecycle
2. **Configuration** (`internal/config`) - Handles configuration loading and validation
3. **Logging** (`internal/logger`) - Provides structured logging capabilities
4. **Raindrop Client** (`internal/raindrop`) - Interacts with the Raindrop.io API
5. **Telegram Bot** (`internal/telegram`) - Handles Telegram bot functionality
6. **Storage** (`internal/storage`) - Manages database operations
7. **Server** (`internal/server`) - Handles OAuth callbacks

## Usage

1. Create a .env file based on `.env.example`:
   ```bash
   cp .env.example .env
   ```

2. Fill in the required environment variables in `.env`:
   - Telegram bot token
   - Raindrop.io OAuth credentials
   - Database connection string

3. Run the application:
   ```bash
   go run cmd/bot/main.go
   ```
   

## Building

To build the application:
```bash
go build -o telegram-raindrop cmd/bot/main.go
```

## Docker

The application includes Docker support:
```bash
docker-compose up
```

## Testing

Run tests with:
```bash
go test ./...
```

## How It Works

1. Users start the bot and are prompted to authenticate with Raindrop.io
2. The bot generates an OAuth link for users to authorize the application
3. After authorization, the OAuth callback is handled by the server component
4. User credentials are stored in the database with proper encryption
5. Users can then send links to the bot, which are saved to their Raindrop.io account
