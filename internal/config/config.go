package config

import (
	"os"
	"strings"

	"github.com/Coiiap5e/TgBotPh/internal/errors"
	"github.com/joho/godotenv"
)

const (
	EnvTelegramBotToken  = "TELEGRAM_BOT_TOKEN"
	EnvTelegramChannelID = "TELEGRAM_CHANNEL_ID"
	EnvKafkaBrokerURLs   = "KAFKA_BROKER_URLS"
	EnvKafkaTopic        = "KAFKA_TOPIC"
	EnvKafkaGroupID      = "KAFKA_GROUP_ID"
)

type Config struct {
	Telegram TelegramConfig
	Kafka    KafkaConfig
}

type TelegramConfig struct {
	BotToken  string
	ChannelID string
}

type KafkaConfig struct {
	BrokerURLs []string
	Topic      string
	GroupID    string
}

func NewConfig() (*Config, error) {

	_ = godotenv.Load()

	// Настройки Telegram
	//telegramBotToken := os.Getenv(EnvTelegramBotToken)
	//if telegramBotToken == "" {
	//	return nil, errors.New(errors.ErrCodeConfig, "Telegram bot token is not set")
	//}
	//telegramChannelID := os.Getenv(EnvTelegramChannelID)
	//if telegramChannelID == "" {
	//	return nil, errors.New(errors.ErrCodeConfig, "Telegram channel ID is not set")
	//}

	// Настройки Kafka
	kafkaBrokerURLsStr := os.Getenv(EnvKafkaBrokerURLs)
	if kafkaBrokerURLsStr == "" {
		return nil, errors.New(errors.ErrCodeConfig, "Kafka broker URLs are not set")
	}
	kafkaBrokerURLs := strings.Split(kafkaBrokerURLsStr, ",")
	if len(kafkaBrokerURLs) == 0 {
		return nil, errors.New(errors.ErrCodeConfig, "No Kafka broker URLs provided")
	}

	kafkaTopic := os.Getenv(EnvKafkaTopic)
	if kafkaTopic == "" {
		return nil, errors.New(errors.ErrCodeConfig, "Kafka topic is not set")
	}

	kafkaGroupID := getEnv(EnvKafkaGroupID, "telegram-bot-consumer-group") // Group ID по умолчанию

	return &Config{
		//Telegram: TelegramConfig{
		//	BotToken:  telegramBotToken,
		//	ChannelID: telegramChannelID,
		//},
		Kafka: KafkaConfig{
			BrokerURLs: kafkaBrokerURLs,
			Topic:      kafkaTopic,
			GroupID:    kafkaGroupID,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
