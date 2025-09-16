package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Messages struct {
	Responses `mapstructure:"responses"`
	Errors    `mapstructure:"errors"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	UnknownCommand    string `mapstructure:"unknown_command"`
	LinkSaved         string `mapstructure:"link_saved"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	UnableToSave string `mapstructure:"unable_to_save"`
	FailAuth     string `mapstructure:"fail_auth"`
}

type Raindrop struct {
	ClientId     string `mapstructure:"client_id" validate:"required"`
	ClientSecret string `mapstructure:"client_secret" validate:"required"`
	RedirectUrl  string `mapstructure:"redirect_url" validate:"required"`
}

type Config struct {
	BotName       string `mapstructure:"bot_name"`
	ServerAddr    string `mapstructure:"server_addr"`
	DatabaseDSN   string `mapstructure:"db_dsn" validate:"required"`
	TelegramToken string `mapstructure:"telegram_token" validate:"required"`
	DebugMode     bool   `mapstructure:"debug_mode"`
	Raindrop      Raindrop
	Messages      Messages
}

func Init() (*Config, error) {
	cfg := Config{}

	if err := loadENV(&cfg); err != nil {
		return nil, fmt.Errorf("loading env failed: %w", err)
	}

	if err := loadMessages(&cfg); err != nil {
		return nil, fmt.Errorf("loading messages failed: %w", err)
	}

	if err := validateCfg(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func loadENV(cfg *Config) error {
	envViper := viper.New()

	envViper.SetConfigFile(".env")
	envViper.AutomaticEnv()
	envViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(envViper)

	if err := envViper.ReadInConfig(); err != nil {
		return err
	}

	if err := envViper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}

func setDefaults(viper *viper.Viper) {
	viper.SetDefault("bot_name", "raindrop_links_bot")
	viper.SetDefault("server_addr", "0.0.0.0:8080")
	viper.SetDefault("debug_mode", true)
}

func loadMessages(cfg *Config) error {
	messagesViper := viper.New()
	messagesViper.AddConfigPath(".")
	messagesViper.SetConfigFile("messages.yml")

	if err := messagesViper.ReadInConfig(); err != nil {
		return err
	}

	if err := messagesViper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}

func validateCfg(cfg *Config) error {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		return err
	}

	return nil
}
