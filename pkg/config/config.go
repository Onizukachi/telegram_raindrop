package config

import (
	"errors"
	"os"
	"strconv"
)

// разбить конфиг по структурами дя базы отдельно рейндроп отдельно
// сделать логер свой и захаркдкодить соообщения в константы
type Config struct {
	BotName       string
	ServerAddr    string
	DatabaseDSN   string
	TelegramToken string
	ClientId      string
	ClientSecret  string
	RedirectUrl   string
	TestBearer    string
	DebugMode     bool
}

func Init() (*Config, error) {
	botName := os.Getenv("BOT_NAME")
	if botName == "" {
		return nil, errors.New("env botName is required")
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		return nil, errors.New("env SERVER_ADDR is required")
	}

	databaseDSN := os.Getenv("DB_DSN")
	if databaseDSN == "" {
		return nil, errors.New("env DB_DSN is required")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return nil, errors.New("env TELEGRAM_TOKEN is required")
	}

	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("env CLIENT_ID is required")
	}

	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("env CLIENT_SECRET is required")
	}

	redirectUrl := os.Getenv("REDIRECT_URL")
	if redirectUrl == "" {
		return nil, errors.New("env REDIRECT_URL is required")
	}

	testBearer := os.Getenv("TEST_BEARER")
	if testBearer == "" {
		return nil, errors.New("env TEST_BEARER is required")
	}

	debugMode := os.Getenv("DEBUG_MODE")
	parseDebugMode, err := strconv.ParseBool(debugMode)
	if err != nil {
		parseDebugMode = true
	}

	cfg := &Config{
		BotName:       botName,
		ServerAddr:    serverAddr,
		DatabaseDSN:   databaseDSN,
		TelegramToken: telegramToken,
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		RedirectUrl:   redirectUrl,
		TestBearer:    testBearer,
		DebugMode:     parseDebugMode,
	}

	return cfg, nil
}
