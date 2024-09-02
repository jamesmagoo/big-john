package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env               string `mapstructure:"ENV"`
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	OpenAIAPIKey      string `mapstructure:"OPENAI_API_KEY"`
	TelegramAuthToken string `mapstructure:"TELEGRAM_AUTH_TOKEN"`
	LogLevel          int    `mapstructure:"LOG_LEVEL"`
}

// reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
