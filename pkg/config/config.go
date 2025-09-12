package config

import (
	"strings"

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
}

type Raindrop struct {
	ClientId     string `mapstructure:"raindrop_client_id"`
	ClientSecret string `mapstructure:"raindrop_client_secret"`
	RedirectUrl  string `mapstructure:"raindrop_redirect_url"`
}

type Config struct {
	BotName       string   `mapstructure:"bot_name"`
	ServerAddr    string   `mapstructure:"server_addr"`
	DatabaseDSN   string   `mapstructure:"db_dsn"`
	TelegramToken string   `mapstructure:"telegram_token"`
	DebugMode     bool     `mapstructure:"debug_mode"`
	Raindrop      Raindrop `mapstructure:",squash"`
	Messages      Messages
}

func Init() (*Config, error) {
	cfg := Config{}

	if err := loadENV(&cfg); err != nil {
		return nil, err
	}

	if err := loadMessages(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadENV(cfg *Config) error {
	envViper := viper.New()
	envViper.SetConfigFile(".env")
	envViper.AutomaticEnv()
	envViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := envViper.ReadInConfig(); err != nil {
		return err
	}

	if err := envViper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
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
