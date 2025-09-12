package config

import (
	"testing"
)

func TestInit(t *testing.T) {
	cfg, err := Init()
	if err != nil {
		t.Fatalf("Error during loading config: %v", err)
	}

	// Test that environment variables are loaded
	if cfg.BotName == "" {
		t.Error("BotName should not be empty")
	}

	if cfg.TelegramToken == "" {
		t.Error("TelegramToken should not be empty")
	}

	if cfg.ServerAddr == "" {
		t.Error("ServerAddr should not be empty")
	}

	if cfg.DatabaseDSN == "" {
		t.Error("DatabaseDSN should not be empty")
	}

	// Test that Raindrop config is loaded
	if cfg.Raindrop.ClientId == "" {
		t.Error("Raindrop.ClientId should not be empty")
	}

	if cfg.Raindrop.ClientSecret == "" {
		t.Error("Raindrop.ClientSecret should not be empty")
	}

	if cfg.Raindrop.RedirectUrl == "" {
		t.Error("Raindrop.RedirectUrl should not be empty")
	}

	// Test that messages are loaded
	if cfg.Messages.Responses.Start == "" {
		t.Error("Messages.Responses.Start should not be empty")
	}

	if cfg.Messages.Responses.AlreadyAuthorized == "" {
		t.Error("Messages.Responses.AlreadyAuthorized should not be empty")
	}

	if cfg.Messages.Responses.UnknownCommand == "" {
		t.Error("Messages.Responses.UnknownCommand should not be empty")
	}

	if cfg.Messages.Responses.LinkSaved == "" {
		t.Error("Messages.Responses.LinkSaved should not be empty")
	}

	if cfg.Messages.Errors.Default == "" {
		t.Error("Messages.Errors.Default should not be empty")
	}

	if cfg.Messages.Errors.InvalidURL == "" {
		t.Error("Messages.Errors.InvalidURL should not be empty")
	}

	if cfg.Messages.Errors.UnableToSave == "" {
		t.Error("Messages.Errors.UnableToSave should not be empty")
	}
}
