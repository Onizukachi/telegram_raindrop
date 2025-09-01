package config

import (
	"errors"
	"os"
)

type Config struct {
	TelegramToken string
	ClientId      string
	ClientSecret  string
	TestBearer    string
}

func Init() (*Config, error) {
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

	testBearer := os.Getenv("TEST_BEARER")
	if testBearer == "" {
		return nil, errors.New("env TEST_BEARER is required")
	}

	cfg := &Config{
		TelegramToken: telegramToken,
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		TestBearer:    testBearer,
	}

	return cfg, nil
}
