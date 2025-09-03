package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DbDSN         string
	TelegramToken string
	ClientId      string
	ClientSecret  string
	RedirectUrl   string
	TestBearer    string
	DebugMode     bool
}

func Init() (*Config, error) {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
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
		DbDSN:         dbDSN,
		TelegramToken: telegramToken,
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		TestBearer:    testBearer,
		DebugMode:     parseDebugMode,
	}

	return cfg, nil
}
