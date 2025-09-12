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
	ClientId     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	RedirectUrl  string `mapstructure:"REDIRECT_URL"`
}

type Config struct {
	BotName       string `mapstructure:"BOT_NAME"`
	ServerAddr    string `mapstructure:"SERVER_ADDR"`
	DatabaseDSN   string `mapstructure:"DB_DSN"`
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
	DebugMode     bool   `mapstructure:"DEBUG_MODE"`
	Raindrop
	Messages
}

func Init() (*Config, error) {
	var cfg Config

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
